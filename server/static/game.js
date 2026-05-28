// game.js — 吉伊卡哇：像素大討伐 Web Client
// cannon-agent + target-system-agent + game-state-agent 負責維護

const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d');

// ── 畫布設定 ──────────────────────────────────────────────────
function resizeCanvas() {
  const game = document.getElementById('game');
  canvas.width = game.clientWidth;
  canvas.height = game.clientHeight - 90; // 扣掉 top/bottom HUD
}
resizeCanvas();

// ── 遊戲狀態 ──────────────────────────────────────────────────
const state = {
  coins: 10000, betLevel: 1, betCost: 1,
  charId: 'chiikawa', charName: 'Chiikawa',
  laborValue: 0, isAuto: false, lockTargetId: '',
  gameState: 'normal_play',
  fireRate: 2.0, projectileSpeed: 700,
};

const targets = new Map();   // instanceId -> target obj
const bullets = [];          // active bullets
const effects = [];          // active effects
const floatTexts = [];       // floating reward texts

let autoFireTimer = 0;
let lastTime = 0;

// ── 目標物顏色 ────────────────────────────────────────────────
const TARGET_COLORS = {
  T001: '#3cb83c', T002: '#50c850', T003: '#e03030',
  T004: '#3060e0', T005: '#ffd040', T006: '#8b5a2b',
  T101: '#888899', T102: '#d4a020', T103: '#fffff0',
  T104: '#ffd700', T105: '#ffcc33', B001: '#cc2244',
};

const CHAR_COLORS = {
  chiikawa: '#ffaacc', hachiware: '#6699ff', usagi: '#ffee66'
};

// ── WebSocket ─────────────────────────────────────────────────
let ws = null;
let reconnectDelay = 1000;
let reconnectTimer = null;

function connect() {
  const proto = location.protocol === 'https:' ? 'wss' : 'ws';
  const host = location.hostname || 'localhost';
  const port = location.port || '7777';
  const playerId = 'player_' + Math.floor(Math.random() * 999999).toString().padStart(6, '0');
  const url = `${proto}://${host}:${port}/ws?player_id=${playerId}`;
  console.log('[WS] Connecting to', url);

  ws = new WebSocket(url);
  ws.onopen = () => {
    console.log('[WS] Connected');
    document.getElementById('disconnect').style.display = 'none';
    reconnectDelay = 1000;
  };
  ws.onclose = () => {
    console.log('[WS] Disconnected');
    document.getElementById('disconnect').style.display = 'flex';
    reconnectTimer = setTimeout(connect, reconnectDelay);
    reconnectDelay = Math.min(reconnectDelay * 2, 30000);
  };
  ws.onerror = (e) => console.error('[WS] Error', e);
  ws.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data);
      handleMessage(msg.type, msg.payload || {});
    } catch(err) { console.error('[WS] Parse error', err); }
  };
}

function send(type, payload) {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type, payload }));
  }
}

// ── 訊息處理 ──────────────────────────────────────────────────
function handleMessage(type, payload) {
  switch(type) {
    case 'game_state':
      state.gameState = payload.state;
      document.getElementById('status').textContent = payload.state.toUpperCase().replace(/_/g,' ');
      break;
    case 'player_update':
      state.coins = payload.coins;
      state.betLevel = payload.bet_level;
      state.betCost = payload.bet_cost;
      state.charId = payload.character_id;
      state.charName = payload.character_name;
      state.laborValue = payload.labor_value;
      state.isAuto = payload.is_auto;
      state.lockTargetId = payload.lock_target_id || '';
      state.fireRate = payload.fire_rate || 2.0;
      state.projectileSpeed = payload.projectile_speed || 700;
      updateHUD();
      break;
    case 'target_spawn':
      spawnTarget(payload);
      break;
    case 'target_update':
      updateTarget(payload);
      break;
    case 'target_kill':
      killTarget(payload);
      break;
    case 'attack_result':
      if (payload.is_hit) {
        spawnHitEffect(payload.target_id);
        screenShake(3);
      }
      break;
    case 'reward':
      showRewardPopup(payload.amount, payload.multiplier);
      break;
    case 'boss_event':
      handleBossEvent(payload);
      break;
    case 'bonus_event':
      handleBonusEvent(payload);
      break;
    case 'announce':
      showAnnounce(payload.message, payload.color || '#ffd700');
      break;
    default:
      // Lucky 系統事件：統一顯示公告橫幅
      if (type.startsWith('lucky_')) {
        handleLuckyEvent(type, payload);
      }
      break;
  }
}

// ── HUD 更新 ──────────────────────────────────────────────────
function updateHUD() {
  document.getElementById('coins').textContent = '💰 ' + state.coins;
  document.getElementById('bet-info').textContent = `BET LV${state.betLevel} (${state.betCost})`;
  document.getElementById('char-name').textContent = state.charName;
  document.getElementById('char-name').style.color = CHAR_COLORS[state.charId] || '#fff';
  const pct = Math.min(100, state.laborValue);
  document.getElementById('labor-bar').style.width = pct + '%';
  document.getElementById('labor-bar').style.background = pct >= 80 ? '#aaaa00' : '#2a8a2a';
  document.getElementById('labor-text').textContent = state.laborValue + '/100';
  document.getElementById('labor-text').style.color = pct >= 80 ? '#ffee00' : '#fff';
  const autoBtn = document.getElementById('auto-btn');
  autoBtn.textContent = state.isAuto ? 'AUTO ON' : 'AUTO';
  autoBtn.className = 'btn' + (state.isAuto ? ' active' : '');
}

// ── 目標物管理 ────────────────────────────────────────────────
function spawnTarget(data) {
  const t = {
    id: data.instance_id,
    defId: data.def_id,
    name: data.name,
    type: data.type,
    x: data.x,
    y: data.y - 40, // 扣掉 top HUD
    hp: data.hp,
    maxHp: data.max_hp,
    speed: data.speed,
    behavior: data.behavior,
    multiplier: data.multiplier,
    isFleeing: false,
    scale: 0, // 進場動畫
    spawnTime: Date.now(),
  };
  targets.set(data.instance_id, t);
}

function updateTarget(data) {
  const t = targets.get(data.instance_id);
  if (!t) return;
  t.hp = data.hp;
  t.x = data.x;
  t.y = data.y - 40;
  if (data.is_fleeing) t.isFleeing = true;
  // 受擊閃白
  t.flashTimer = 0.12;
}

function killTarget(data) {
  const t = targets.get(data.instance_id);
  if (t) {
    spawnKillEffect(t.x, t.y, data.multiplier);
    if (data.reward > 0) {
      addFloatText(t.x, t.y, '+' + data.reward + ' x' + data.multiplier.toFixed(0), data.multiplier);
    }
    targets.delete(data.instance_id);
  }
}

// ── 射擊 ──────────────────────────────────────────────────────
const CANNON_X = 640;
const CANNON_Y = canvas.height - 30;

canvas.addEventListener('click', (e) => {
  if (!['normal_play','boss_battle'].includes(state.gameState)) return;
  const rect = canvas.getBoundingClientRect();
  const cx = e.clientX - rect.left;
  const cy = e.clientY - rect.top;
  fireAt(cx, cy);
});

function fireAt(cx, cy) {
  // 找最近的目標
  let bestId = '';
  let bestDist = 80;
  for (const [id, t] of targets) {
    const d = Math.hypot(t.x - cx, t.y - cy);
    if (d < bestDist) { bestDist = d; bestId = id; }
  }
  send('attack', { target_id: bestId, click_x: cx, click_y: cy });
  spawnBullet(cx, cy, bestId);
}

function spawnBullet(tx, ty, targetId) {
  const color = CHAR_COLORS[state.charId] || '#fff';
  const dist = Math.hypot(tx - CANNON_X, ty - CANNON_Y);
  const speed = state.projectileSpeed || 700;
  const time = Math.max(0.05, Math.min(0.3, dist / speed));
  bullets.push({ x: CANNON_X, y: CANNON_Y, tx, ty, color, time, elapsed: 0, targetId });
}

// ── AUTO 射擊 ─────────────────────────────────────────────────
function autoFire(delta) {
  if (!state.isAuto) return;
  if (!['normal_play','boss_battle'].includes(state.gameState)) return;
  autoFireTimer += delta;
  const interval = 1.0 / (state.fireRate || 2.0);
  if (autoFireTimer < interval) return;
  autoFireTimer = 0;

  // 找最高價值目標
  let best = null, bestScore = -1;
  for (const [, t] of targets) {
    let score = t.multiplier * 2;
    if (t.x < 400) score += 20;
    if (t.type === 'boss') score += 500;
    if (score > bestScore) { bestScore = score; best = t; }
  }
  if (!best) return;
  send('attack', { target_id: best.id, click_x: best.x, click_y: best.y });
  spawnBullet(best.x, best.y, best.id);
}

// ── 特效 ──────────────────────────────────────────────────────
function spawnHitEffect(targetId) {
  const t = targets.get(targetId);
  if (!t) return;
  effects.push({ type: 'hit', x: t.x, y: t.y, r: 8, maxR: 40, life: 0.15, elapsed: 0, color: CHAR_COLORS[state.charId] || '#fff' });
}

function spawnKillEffect(x, y, mult) {
  const color = mult >= 50 ? '#ff6600' : mult >= 20 ? '#ffd700' : '#aaaaff';
  effects.push({ type: 'kill', x, y, r: 10, maxR: 60, life: 0.25, elapsed: 0, color });
  // 粒子
  for (let i = 0; i < 6; i++) {
    const angle = (i / 6) * Math.PI * 2;
    effects.push({ type: 'particle', x, y, vx: Math.cos(angle) * 80, vy: Math.sin(angle) * 80, life: 0.4, elapsed: 0, color });
  }
}

// ── 浮動文字 ──────────────────────────────────────────────────
function addFloatText(x, y, text, mult) {
  const color = mult >= 100 ? '#ff4400' : mult >= 20 ? '#ffd700' : '#ffffff';
  floatTexts.push({ x, y, text, color, life: 0.8, elapsed: 0 });
}

// ── 震動 ──────────────────────────────────────────────────────
let shakeX = 0, shakeY = 0, shakeMag = 0;
function screenShake(mag) { shakeMag = Math.max(shakeMag, mag); }

// ── 公告 ──────────────────────────────────────────────────────
let announceTimer = null;
function showAnnounce(msg, color) {
  const el = document.getElementById('announce');
  el.textContent = msg;
  el.style.color = color || '#ffd700';
  el.style.opacity = '1';
  if (announceTimer) clearTimeout(announceTimer);
  announceTimer = setTimeout(() => { el.style.opacity = '0'; }, 3000);
}

// ── 獎勵彈窗 ──────────────────────────────────────────────────
let rewardTimer = null;
function showRewardPopup(amount, mult) {
  if (amount <= 0) return;
  const el = document.getElementById('reward-popup');
  const icon = mult >= 100 ? '🌟' : mult >= 20 ? '⭐' : '💰';
  el.textContent = `${icon} +${amount}  x${mult.toFixed(0)}`;
  el.style.color = mult >= 100 ? '#ff4400' : mult >= 20 ? '#ffd700' : '#ffffff';
  el.style.opacity = '1';
  el.style.top = '300px';
  if (rewardTimer) clearTimeout(rewardTimer);
  rewardTimer = setTimeout(() => { el.style.opacity = '0'; }, 1500);
}

// ── BOSS / Bonus 事件 ─────────────────────────────────────────
function handleBossEvent(payload) {
  switch(payload.event) {
    case 'warning': showAnnounce('⚠️ BOSS WARNING!', '#ff4444'); break;
    case 'spawn':   showAnnounce('⚔️ BOSS APPEARED!', '#ff2222'); break;
    case 'phase_change': showAnnounce('💢 PHASE 2!', '#ff0000'); break;
    case 'kill':    showAnnounce(`🏆 BOSS KILLED! x${payload.multiplier}`, '#ffd700'); break;
    case 'timeout': showAnnounce('💨 BOSS ESCAPED', '#888888'); break;
  }
}

function handleBonusEvent(payload) {
  switch(payload.event) {
    case 'start':  showAnnounce('🌿 BONUS GAME START!', '#44ff44'); break;
    case 'result': showAnnounce(`🎉 BONUS! x${payload.multiplier.toFixed(1)} +${payload.reward}`, '#ffd700'); break;
  }
}

// ── Lucky 系統事件處理 ────────────────────────────────────────
const LUCKY_NAMES = {
  lucky_chain_lightning: '⚡ 連鎖閃電',
  lucky_crab_torpedo: '🦀 螃蟹魚雷',
  lucky_vortex: '🌀 渦旋海葵',
  lucky_golden_dragon: '🐉 黃金龍魚',
  lucky_thunder_lobster: '🦞 雷霆龍蝦',
  lucky_awakened_phoenix: '🔥 覺醒鳳凰',
  lucky_shockwave_bomb: '💥 全場震盪',
  lucky_drill_torpedo: '🚀 鑽頭魚雷',
  lucky_time_freeze: '❄️ 時間凍結',
  lucky_chain_explosion: '💥 連鎖爆炸',
  lucky_chain_long_king: '👑 千龍王輪盤',
  lucky_dragon_shotgun: '🐲 龍力散彈',
  lucky_rocket_cannon: '🚀 火箭砲',
  lucky_deep_whirlpool: '🌊 深海漩渦',
  lucky_vampire_mult: '🧛 吸血鬼',
  lucky_jackpot_pool: '🎰 Progressive Jackpot',
  lucky_dragon_king: '👑 龍王輪盤',
  lucky_genesis_epoch: '🌌 創世紀元',
  lucky_energy_storm: '⚡ 能量風暴',
  lucky_crystal_resonance: '💎 水晶共鳴',
  lucky_fate_judgment: '⚖️ 命運審判',
  lucky_time_reversal: '⏪ 時間逆流',
  lucky_cosmic_singularity: '🌌 宇宙奇點',
};

function handleLuckyEvent(type, payload) {
  const event = payload.event || '';
  const name = LUCKY_NAMES[type] || type.replace('lucky_', '').replace(/_/g, ' ').toUpperCase();
  const triggerName = payload.trigger_name || payload.player_name || '玩家';
  
  // 只顯示觸發事件（不顯示每個 tick 更新）
  if (event === 'trigger' || event === 'jackpot_win' || event === 'start' || event === 'win') {
    let msg = `✨ ${name}`;
    if (triggerName) msg += ` — ${triggerName}`;
    if (payload.multiplier) msg += ` ×${payload.multiplier.toFixed ? payload.multiplier.toFixed(1) : payload.multiplier}`;
    if (payload.reward) msg += ` +${payload.reward}`;
    if (payload.tier_name) msg += ` [${payload.tier_name}]`;
    showAnnounce(msg, '#ffd700');
    screenShake(5);
  } else if (event === 'settle' || event === 'end') {
    let msg = `🏆 ${name} 結算`;
    if (payload.total_reward) msg += ` +${payload.total_reward}`;
    showAnnounce(msg, '#ffaa00');
  }
}

// ── 控制函數 ──────────────────────────────────────────────────
function toggleAuto() { send('auto_toggle', {}); }
function sendBetChange(dir) {
  const newLevel = Math.max(1, Math.min(10, state.betLevel + dir));
  send('bet_change', { bet_level: newLevel });
}
function sendLock(id) { send('lock', { target_id: id }); }

// ── 主循環 ────────────────────────────────────────────────────
function gameLoop(timestamp) {
  const delta = Math.min(0.05, (timestamp - lastTime) / 1000);
  lastTime = timestamp;

  // 震動衰減
  shakeMag = Math.max(0, shakeMag - 20 * delta);
  shakeX = (Math.random() - 0.5) * shakeMag;
  shakeY = (Math.random() - 0.5) * shakeMag;

  // AUTO
  autoFire(delta);

  // 更新目標物位置
  for (const [id, t] of targets) {
    // 進場縮放
    if (t.scale < 1) t.scale = Math.min(1, t.scale + delta * 6);
    // 閃白計時
    if (t.flashTimer > 0) t.flashTimer -= delta;
    // 移動
    const speed = t.isFleeing ? t.speed * 2.5 : t.speed;
    switch(t.behavior) {
      case 'linear': case 'flee': case 'fast': t.x -= speed * delta; break;
      case 'sink': t.y += speed * 0.3 * delta; t.x -= 10 * delta; break;
    }
    if (t.x < -100) targets.delete(id);
  }

  // 更新子彈
  for (let i = bullets.length - 1; i >= 0; i--) {
    const b = bullets[i];
    b.elapsed += delta;
    const p = Math.min(1, b.elapsed / b.time);
    b.x = CANNON_X + (b.tx - CANNON_X) * p;
    b.y = CANNON_Y + (b.ty - CANNON_Y) * p;
    if (p >= 1) bullets.splice(i, 1);
  }

  // 更新特效
  for (let i = effects.length - 1; i >= 0; i--) {
    const e = effects[i];
    e.elapsed += delta;
    if (e.type === 'particle') { e.x += e.vx * delta; e.y += e.vy * delta; }
    if (e.elapsed >= e.life) effects.splice(i, 1);
  }

  // 更新浮動文字
  for (let i = floatTexts.length - 1; i >= 0; i--) {
    const f = floatTexts[i];
    f.elapsed += delta;
    f.y -= 60 * delta;
    if (f.elapsed >= f.life) floatTexts.splice(i, 1);
  }

  draw();
  requestAnimationFrame(gameLoop);
}

// ── 繪製 ──────────────────────────────────────────────────────
function draw() {
  ctx.save();
  ctx.translate(shakeX, shakeY);

  // 清除
  ctx.clearRect(-10, -10, canvas.width + 20, canvas.height + 20);

  // 背景
  const grad = ctx.createLinearGradient(0, 0, 0, canvas.height);
  grad.addColorStop(0, '#0a1628');
  grad.addColorStop(0.5, '#0d2040');
  grad.addColorStop(1, '#0a1830');
  ctx.fillStyle = grad;
  ctx.fillRect(0, 0, canvas.width, canvas.height);

  // 氣泡（裝飾）
  drawBubbles();

  // 目標物
  for (const [, t] of targets) drawTarget(t);

  // 子彈
  for (const b of bullets) drawBullet(b);

  // 特效
  for (const e of effects) drawEffect(e);

  // 浮動文字
  for (const f of floatTexts) drawFloatText(f);

  // 砲台角色
  drawCannon();

  ctx.restore();
}

// ── 氣泡 ──────────────────────────────────────────────────────
const bubbles = Array.from({length: 20}, () => ({
  x: Math.random() * 1280, y: Math.random() * 630 + 40,
  r: 3 + Math.random() * 8, speed: 15 + Math.random() * 30,
  alpha: 0.2 + Math.random() * 0.4
}));
function drawBubbles() {
  for (const b of bubbles) {
    b.y -= b.speed * 0.016;
    if (b.y < 40) { b.y = canvas.height; b.x = Math.random() * 1280; }
    ctx.beginPath();
    ctx.arc(b.x, b.y, b.r, 0, Math.PI * 2);
    ctx.strokeStyle = `rgba(150,200,255,${b.alpha})`;
    ctx.lineWidth = 1;
    ctx.stroke();
  }
}

// ── 目標物繪製 ────────────────────────────────────────────────
function drawTarget(t) {
  const s = t.scale;
  const size = t.type === 'boss' ? 80 : 40;
  const color = TARGET_COLORS[t.defId] || '#888';

  ctx.save();
  ctx.translate(t.x, t.y);
  ctx.scale(s, s);

  // 高倍率光暈
  if (t.multiplier >= 30) {
    const glowColor = t.multiplier >= 50 ? 'rgba(255,100,0,0.3)' : 'rgba(255,215,0,0.25)';
    ctx.beginPath();
    ctx.arc(0, 0, size * 0.7, 0, Math.PI * 2);
    ctx.fillStyle = glowColor;
    ctx.fill();
  }

  // 主體
  ctx.beginPath();
  ctx.arc(0, 0, size * 0.5, 0, Math.PI * 2);
  const flash = t.flashTimer > 0 ? Math.min(1, t.flashTimer / 0.12) : 0;
  ctx.fillStyle = flash > 0 ? `rgba(255,255,255,${flash})` : color;
  ctx.fill();
  ctx.strokeStyle = 'rgba(0,0,0,0.6)';
  ctx.lineWidth = 2;
  ctx.stroke();

  // 眼睛
  if (t.type !== 'boss') {
    ctx.fillStyle = '#000';
    ctx.beginPath(); ctx.arc(-7, -4, 3, 0, Math.PI * 2); ctx.fill();
    ctx.beginPath(); ctx.arc(7, -4, 3, 0, Math.PI * 2); ctx.fill();
    ctx.fillStyle = '#fff';
    ctx.beginPath(); ctx.arc(-6, -5, 1.5, 0, Math.PI * 2); ctx.fill();
    ctx.beginPath(); ctx.arc(8, -5, 1.5, 0, Math.PI * 2); ctx.fill();
  } else {
    // BOSS 大眼
    ctx.fillStyle = '#ff2222';
    ctx.beginPath(); ctx.arc(-18, -8, 10, 0, Math.PI * 2); ctx.fill();
    ctx.beginPath(); ctx.arc(18, -8, 10, 0, Math.PI * 2); ctx.fill();
    ctx.fillStyle = '#000';
    ctx.beginPath(); ctx.arc(-18, -8, 5, 0, Math.PI * 2); ctx.fill();
    ctx.beginPath(); ctx.arc(18, -8, 5, 0, Math.PI * 2); ctx.fill();
  }

  // 倍率標籤
  ctx.fillStyle = t.multiplier >= 30 ? '#ffd700' : '#ffffff';
  ctx.font = `bold ${t.type === 'boss' ? 14 : 11}px monospace`;
  ctx.textAlign = 'center';
  ctx.fillText('x' + t.multiplier.toFixed(0), 0, size * 0.5 + 14);

  ctx.restore();

  // HP 條
  const barW = size * 1.2;
  const barH = 5;
  const barX = t.x - barW / 2;
  const barY = t.y - size * 0.5 * s - 12;
  const pct = t.hp / t.maxHp;
  ctx.fillStyle = '#222';
  ctx.fillRect(barX, barY, barW, barH);
  ctx.fillStyle = pct > 0.6 ? '#22cc22' : pct > 0.3 ? '#cccc22' : '#cc2222';
  ctx.fillRect(barX, barY, barW * pct, barH);
}

// ── 子彈繪製 ──────────────────────────────────────────────────
function drawBullet(b) {
  ctx.save();
  ctx.translate(b.x, b.y);
  const angle = Math.atan2(b.ty - CANNON_Y, b.tx - CANNON_X);
  ctx.rotate(angle);
  ctx.fillStyle = b.color;
  ctx.fillRect(-8, -4, 16, 8);
  ctx.fillStyle = 'rgba(255,255,255,0.8)';
  ctx.fillRect(4, -2, 4, 4);
  ctx.restore();
}

// ── 特效繪製 ──────────────────────────────────────────────────
function drawEffect(e) {
  const p = e.elapsed / e.life;
  const alpha = 1 - p;
  ctx.globalAlpha = alpha;
  if (e.type === 'hit' || e.type === 'kill') {
    const r = e.r + (e.maxR - e.r) * p;
    ctx.beginPath();
    ctx.arc(e.x, e.y, r, 0, Math.PI * 2);
    ctx.strokeStyle = e.color;
    ctx.lineWidth = 3;
    ctx.stroke();
  } else if (e.type === 'particle') {
    ctx.fillStyle = e.color;
    ctx.fillRect(e.x - 3, e.y - 3, 6, 6);
  }
  ctx.globalAlpha = 1;
}

// ── 浮動文字繪製 ──────────────────────────────────────────────
function drawFloatText(f) {
  const alpha = 1 - f.elapsed / f.life;
  ctx.globalAlpha = alpha;
  ctx.fillStyle = f.color;
  ctx.font = 'bold 18px monospace';
  ctx.textAlign = 'center';
  ctx.fillText(f.text, f.x, f.y);
  ctx.globalAlpha = 1;
}

// ── 砲台角色繪製 ──────────────────────────────────────────────
function drawCannon() {
  const cx = CANNON_X, cy = canvas.height - 30;
  const color = CHAR_COLORS[state.charId] || '#ffaacc';

  // 身體
  ctx.fillStyle = '#fffff7';
  ctx.strokeStyle = '#292a2b';
  ctx.lineWidth = 2;
  ctx.beginPath();
  ctx.arc(cx, cy - 20, 22, 0, Math.PI * 2);
  ctx.fill(); ctx.stroke();

  // 耳朵
  ctx.beginPath();
  ctx.arc(cx - 16, cy - 38, 8, 0, Math.PI * 2);
  ctx.fill(); ctx.stroke();
  ctx.beginPath();
  ctx.arc(cx + 16, cy - 38, 8, 0, Math.PI * 2);
  ctx.fill(); ctx.stroke();

  // 眼睛
  ctx.fillStyle = '#000';
  ctx.beginPath(); ctx.arc(cx - 7, cy - 22, 3, 0, Math.PI * 2); ctx.fill();
  ctx.beginPath(); ctx.arc(cx + 7, cy - 22, 3, 0, Math.PI * 2); ctx.fill();
  ctx.fillStyle = '#fff';
  ctx.beginPath(); ctx.arc(cx - 6, cy - 23, 1.5, 0, Math.PI * 2); ctx.fill();
  ctx.beginPath(); ctx.arc(cx + 8, cy - 23, 1.5, 0, Math.PI * 2); ctx.fill();

  // 腮紅
  ctx.fillStyle = 'rgba(239,165,201,0.6)';
  ctx.beginPath(); ctx.arc(cx - 13, cy - 16, 5, 0, Math.PI * 2); ctx.fill();
  ctx.beginPath(); ctx.arc(cx + 13, cy - 16, 5, 0, Math.PI * 2); ctx.fill();

  // 角色名
  ctx.fillStyle = color;
  ctx.font = '12px monospace';
  ctx.textAlign = 'center';
  ctx.fillText(state.charName, cx, cy + 8);
}

// ── 啟動 ──────────────────────────────────────────────────────
connect();
requestAnimationFrame(gameLoop);
