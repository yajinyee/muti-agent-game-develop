## GameManager.gd
## 遊戲狀態管理，處理 Server 訊息並更新遊戲
## Autoload 單例

extends Node
class_name GameManagerClass

# 訊號
signal game_state_changed(new_state: String)
signal target_spawned(target_data: Dictionary)
signal target_updated(update_data: Dictionary)
signal target_killed(kill_data: Dictionary)
signal target_teleported(target_id: String, new_pos: Vector2)  # 傳送魚位置同步（DAY-223）
signal freeze_world_started(speed_factor: float)               # 冰凍世界開始（DAY-237）
signal freeze_world_ended()                                    # 冰凍世界結束（DAY-237）
signal attack_result(result: Dictionary)
signal reward_received(reward: Dictionary)
signal player_updated(player_data: Dictionary)
signal boss_event(event_data: Dictionary)
signal bonus_event(event_data: Dictionary)
signal leaderboard_updated(entries: Array)
signal achievement_unlocked(achievement_data: Dictionary)
signal combo_event(combo_data: Dictionary)  # 連擊事件（DAY-022）
signal mission_updated(missions: Array)     # 任務進度更新（DAY-037）
signal mission_completed(mission_data: Dictionary)  # 任務完成（DAY-037）
signal jackpot_updated(jackpot_data: Dictionary)    # Jackpot 池更新（DAY-048）
signal jackpot_won(win_data: Dictionary)            # Jackpot 中獎（DAY-048）
signal jackpot_animation(anim_data: Dictionary)     # Jackpot 觸發動畫（DAY-095）
signal player_stats_updated(stats_data: Dictionary) # 玩家統計更新（DAY-096）
signal announcement_received(ann_data: Dictionary)  # 全服公告（DAY-097）
signal spectator_joined(spectator_data: Dictionary) # 觀戰者加入通知（DAY-054d）
signal daily_bonus_received(bonus_data: Dictionary) # 每日登入獎勵（DAY-065）
signal spectator_left(spectator_data: Dictionary)  # 觀戰者離開通知（DAY-055）
signal tournament_updated(tournament_data: Dictionary) # 週賽排名更新（DAY-066）
signal daily_tournament_updated(tournament_data: Dictionary) # 每日賽排名更新（DAY-093）
signal multi_format_updated(tournament_data: Dictionary)    # 多格式賽排名更新（DAY-111）
signal activity_feed_event(event_data: Dictionary)          # 動態牆新事件（DAY-112）
signal activity_feed_history(history_data: Dictionary)      # 動態牆歷史（DAY-112）
signal roulette_started(roulette_data: Dictionary)          # 雙層輪盤開始（DAY-113）
signal roulette_result(result_data: Dictionary)             # 雙層輪盤結果（DAY-113）
signal buy_bonus_status(status_data: Dictionary)            # Buy Bonus 狀態（DAY-114）
signal buy_bonus_success(success_data: Dictionary)          # Buy Bonus 購買成功（DAY-114）
signal buy_bonus_error(error_data: Dictionary)              # Buy Bonus 購買失敗（DAY-114）
# Co-op Boss Raid 系統（DAY-115）
signal raid_warning(raid_data: Dictionary)                  # 討伐警告廣播
signal raid_started(raid_data: Dictionary)                  # 討伐開始廣播
signal raid_updated(raid_data: Dictionary)                  # 討伐狀態更新
signal raid_result(result_data: Dictionary)                 # 討伐結算廣播
# 碎片收集大獎系統（DAY-116）
signal fragment_dropped(drop_data: Dictionary)              # 碎片掉落通知
signal fragment_completed(complete_data: Dictionary)        # 集齊碎片大獎廣播
signal fragment_status_received(status_data: Dictionary)    # 碎片狀態回應
# 幸運捕獲系統（DAY-119）
signal lucky_catch(catch_data: Dictionary)                  # 幸運捕獲廣播（全服）
# 任務連續寬限期（DAY-120）
signal mission_mercy_protected(mercy_data: Dictionary)      # 寬限期保護通知
# Rapid Respin 系統（DAY-121）
signal rapid_respin(respin_data: Dictionary)                # Rapid Respin 觸發廣播（全服）
signal rapid_respin_end(end_data: Dictionary)               # Rapid Respin 連鎖結束通知
# 寶藏地圖系統（DAY-122）
signal treasure_map_updated(map_data: Dictionary)           # 寶藏地圖狀態更新
signal treasure_map_line(line_data: Dictionary)             # 完成一行/列/對角線通知
signal treasure_map_full(full_data: Dictionary)             # 完成整張地圖通知
# 閃電挑戰系統（DAY-123）
signal flash_challenge_started(data: Dictionary)            # 閃電挑戰開始廣播
signal flash_challenge_updated(data: Dictionary)            # 閃電挑戰進度更新
signal flash_challenge_ended(data: Dictionary)              # 閃電挑戰結束廣播
signal flash_challenge_reward(data: Dictionary)             # 閃電挑戰獎勵通知（個人）
# 傳說目標警報系統（DAY-124）
signal rare_target_alerted(data: Dictionary)                # 稀有/傳說目標出現廣播
# 黃金時間系統（DAY-125）
signal golden_time_started(data: Dictionary)                # 黃金時間開始廣播
signal golden_time_ended(data: Dictionary)                  # 黃金時間結束廣播
signal golden_time_status(data: Dictionary)                 # 黃金時間狀態回應（個人）
# 稀有連擊累積倍率系統（DAY-126）
signal rare_catch_updated(data: Dictionary)                 # 稀有連擊更新（個人）
signal rare_catch_broadcasted(data: Dictionary)             # 稀有連擊廣播（全服）
signal rare_catch_reset(data: Dictionary)                   # 稀有連擊重置（個人）

# 天氣湧現事件（DAY-127）
signal weather_surge_started(data: Dictionary)              # 天氣湧現開始（全服廣播）
signal weather_surge_ended(data: Dictionary)                # 天氣湧現結束（全服廣播）

# 龍怒蓄力大招系統（DAY-128）
signal wrath_updated(data: Dictionary)                      # 怒氣值更新（個人）
signal wrath_started(data: Dictionary)                      # 大招開始（全服廣播）
signal wrath_result(data: Dictionary)                       # 大招結果（全服廣播）
# 不死 BOSS 連勝系統（DAY-129）
signal immortal_boss_spawned(data: Dictionary)              # 不死 BOSS 出現（全服廣播）
signal immortal_boss_hit(data: Dictionary)                  # 命中不死 BOSS（全服廣播）
signal immortal_boss_left(data: Dictionary)                 # 不死 BOSS 離開（全服廣播）
signal immortal_boss_status(data: Dictionary)               # 不死 BOSS 狀態（個人）
# 覺醒 BOSS 系統（DAY-130）
signal awaken_boss_spawned(data: Dictionary)                # 覺醒 BOSS 出現（全服廣播）
signal awaken_boss_hit(data: Dictionary)                    # 命中覺醒 BOSS（全服廣播）
signal awaken_boss_powerup(data: Dictionary)                # Power Up 觸發（全服廣播）
signal awaken_boss_left(data: Dictionary)                   # 覺醒 BOSS 離開（全服廣播）
signal awaken_boss_status(data: Dictionary)                 # 覺醒 BOSS 狀態（個人）
# 連勝獎勵系統（DAY-131）
signal win_streak_updated(data: Dictionary)                 # 連勝更新（個人）
signal win_streak_milestone(data: Dictionary)               # 里程碑達成（個人/全服）
signal win_streak_reset(data: Dictionary)                   # 連勝重置（個人）
# 閃電鰻連鎖攻擊系統（DAY-132）
signal lightning_eel_chain(data: Dictionary)                # 連鎖攻擊結果廣播（全服）
signal lightning_eel_status(data: Dictionary)               # 閃電鰻冷卻狀態（個人）
# 狂熱模式系統（DAY-133）
signal fever_mode_started(data: Dictionary)                 # 狂熱模式開始（全服廣播）
signal fever_mode_ended(data: Dictionary)                   # 狂熱模式結束（個人）
signal fever_mode_status(data: Dictionary)                  # 狂熱模式狀態更新（個人）
signal title_unlocked(title_data: Dictionary)          # 稱號解鎖通知（DAY-068）
signal skin_updated(skin_data: Dictionary)             # 砲台外觀更新（DAY-071）
signal season_updated(season_data: Dictionary)         # 賽季通行證更新（DAY-072）
signal season_level_up(level_data: Dictionary)         # 賽季等級升級（DAY-072）
signal friend_list_updated(friend_data: Dictionary)    # 好友列表更新（DAY-073）
signal friend_request_received(request_data: Dictionary) # 好友請求通知（DAY-073）
signal friend_updated(update_data: Dictionary)         # 好友狀態更新（DAY-073）
# 好友禮物系統（DAY-101）
signal gift_received(gift_data: Dictionary)            # 收到禮物通知
signal gift_sent(gift_data: Dictionary)                # 送出禮物成功通知
signal gift_status(status_data: Dictionary)            # 今日禮物狀態
signal gift_error(error_data: Dictionary)              # 禮物操作失敗
# 好友挑戰系統（DAY-102）
signal challenge_request(request_data: Dictionary)     # 收到挑戰邀請
signal challenge_updated(update_data: Dictionary)      # 挑戰狀態/分數更新
signal challenge_result(result_data: Dictionary)       # 挑戰結果
signal challenge_error(error_data: Dictionary)         # 挑戰操作失敗
# 私訊系統（DAY-103）
signal dm_received(dm_data: Dictionary)                # 收到私訊
signal dm_sent(dm_data: Dictionary)                    # 發送成功確認
signal dm_error(error_data: Dictionary)                # 發送失敗
signal open_dm_panel(friend_id: String, friend_name: String) # 開啟 DM 面板
# 玩家名片系統（DAY-106）
signal player_card_received(card_data: Dictionary)     # 收到玩家名片資料
# 登入里程碑系統（DAY-107）
signal login_milestone_reached(milestone_data: Dictionary) # 里程碑達成通知
signal login_progress_received(progress_data: Dictionary)  # 登入進度回應
# 超級 Bonus 系統（DAY-108）
signal super_bonus_triggered(bonus_data: Dictionary)       # 超級 Bonus 觸發通知
signal guild_updated(guild_data: Dictionary)           # 公會資訊更新（DAY-074）
signal guild_task_complete(task_data: Dictionary)      # 公會任務完成（DAY-074）
signal guild_message_received(msg_data: Dictionary)    # 公會聊天訊息（DAY-075）
signal guild_war_updated(war_data: Dictionary)         # 公會戰排名更新（DAY-076）
signal guild_war_result(result_data: Dictionary)       # 公會戰結算通知（DAY-076）
signal daily_boss_updated(boss_data: Dictionary)       # 每日 BOSS 狀態更新（DAY-077）
signal daily_boss_defeated(defeat_data: Dictionary)    # 每日 BOSS 擊殺通知（DAY-077）
signal vip_updated(vip_data: Dictionary)               # VIP 狀態更新（DAY-078）
signal vip_level_up(level_data: Dictionary)            # VIP 升級通知（DAY-078）
signal vip_weekly_claimed(claim_data: Dictionary)      # VIP 週獎勵領取通知（DAY-078）
signal event_updated(event_data: Dictionary)           # 限時活動狀態更新（DAY-079）
signal codex_updated(codex_data: Dictionary)           # 圖鑑狀態更新（DAY-081）
signal codex_unlocked(unlock_data: Dictionary)         # 圖鑑條目解鎖通知（DAY-081）
signal codex_complete(complete_data: Dictionary)       # 全圖鑑完成通知（DAY-081）
signal streak_updated(streak_data: Dictionary)         # 連擊狀態更新（DAY-083）
signal streak_reset(reset_data: Dictionary)            # 連擊重置通知（DAY-083）
signal referral_info_received(info_data: Dictionary)   # 推薦碼資訊（DAY-082）
signal referral_success(success_data: Dictionary)      # 推薦碼使用成功（DAY-082）
signal referral_error(error_data: Dictionary)          # 推薦碼使用失敗（DAY-082）
signal wheel_triggered(wheel_data: Dictionary)         # 幸運轉盤觸發（DAY-084）
signal challenge_unlocked(challenge_data: Dictionary)  # 隱藏挑戰解鎖（DAY-085）
signal mission_streak_bonus(streak_data: Dictionary)   # 任務連續完成獎勵（DAY-086）
signal weather_updated(weather_data: Dictionary)       # 天氣狀態更新（DAY-087）
signal chain_explosion(chain_data: Dictionary)         # 連鎖爆炸（DAY-088）
signal chain_target_killed(instance_id: String, multiplier: float) # 連鎖目標擊破（DAY-088）
signal special_weapon_updated(weapon_data: Dictionary) # 特殊武器狀態更新（DAY-089）
signal special_weapon_fired(fire_data: Dictionary)     # 特殊武器發射廣播（DAY-089）
signal special_weapon_charged(charge_data: Dictionary) # 特殊武器自動充能完成（DAY-134）
signal homing_missile_result(result_data: Dictionary)  # 追蹤飛彈命中結果（DAY-141）
signal dragon_wrath_charge(charge_data: Dictionary)    # 龍怒流星雨怒氣值更新（DAY-154）
signal dragon_wrath_result(result_data: Dictionary)    # 龍怒流星雨結果（DAY-154）
signal torpedo_result(result_data: Dictionary)         # 魚雷爆炸結果（DAY-155）
signal railgun_result(result_data: Dictionary)         # 軌道炮穿透結果（DAY-157）
signal black_hole_result(result_data: Dictionary)      # 黑洞漩渦爆炸結果（DAY-166）
signal roulette_crab_start(data: Dictionary)           # 黃金輪盤螃蟹開始（DAY-167）
signal roulette_crab_result(data: Dictionary)          # 黃金輪盤螃蟹結果（DAY-167）
signal lion_dance_burst(data: Dictionary)              # 獅子舞大獎爆發（DAY-168）
signal vortex_fish(data: Dictionary)                   # 漩渦魚群吸引（DAY-169）
signal freeze_bomb(data: Dictionary)                   # 冰凍炸彈魚（DAY-170）
signal ice_fishing_wheel(data: Dictionary)             # 冰釣幸運輪盤（DAY-171）
signal lucky_egg_fish(data: Dictionary)                # 幸運彩蛋魚（DAY-172）
signal rainbow_lucky_fish(data: Dictionary)            # 彩虹幸運魚（DAY-173）
signal sea_anemone(data: Dictionary)                   # 海葵觸手攻擊（DAY-174）
signal lucky_dice_fish(data: Dictionary)               # 幸運骰子魚（DAY-175）
signal fire_storm_fish(data: Dictionary)               # 火焰風暴魚（DAY-176）
signal golden_treasure_fish(data: Dictionary)          # 黃金寶藏魚（DAY-177）
signal mermaid_healing(data: Dictionary)               # 美人魚治癒（DAY-178）
signal lucky_clover_fish(data: Dictionary)             # 幸運草魚（DAY-179）
signal rainbow_shark_burst(data: Dictionary)           # 彩虹鯊魚爆發（DAY-180）
signal thunder_shark_chain(data: Dictionary)           # 雷霆鯊魚連鎖閃電（DAY-181）
signal vampire_fish(data: Dictionary)                  # 吸血鬼魚累積倍率（DAY-182）
signal lightning_auto_chain(data: Dictionary)          # 閃電魚自動連鎖（DAY-183）
signal meteor_fish(data: Dictionary)                   # 隕石魚隕石雨（DAY-184）
signal phoenix_fish(data: Dictionary)                  # 鳳凰魚涅槃重生（DAY-185）
signal dragon_turtle(data: Dictionary)                 # 龍龜不死 Boss（DAY-186）
signal chain_bomb(data: Dictionary)                    # 連鎖爆炸魚（DAY-187）
signal crocodile_hunter(data: Dictionary)              # 巨型鱷魚獵食（DAY-188）
signal time_bomb_fish(data: Dictionary)                # 時間炸彈魚（DAY-189）
signal triple_lucky_fish(data: Dictionary)             # 三重幸運魚（DAY-190）
signal school_panic(data: Dictionary)                  # 魚群驚嚇連帶（DAY-191）
signal rock_skeleton_concert(data: Dictionary)         # 搖滾骷髏演唱會（DAY-192）
signal electric_jellyfish(data: Dictionary)            # 電流水母電流網路（DAY-193）
signal chainlong_king(data: Dictionary)                # 長龍王雙環輪盤（DAY-194）
signal drill_lobster(data: Dictionary)                 # 鑽頭龍蝦穿透爆炸（DAY-195）
signal anglerfish_electric(data: Dictionary)           # 巨型鮟鱇魚電擊寶箱（DAY-196）
signal mystic_dragon(data: Dictionary)                 # 神秘龍魚八波攻擊（DAY-197）
signal ghost_fish(data: Dictionary)                    # 幽靈魚分身（DAY-198）
signal thunderbolt_lobster(data: Dictionary)           # 雷霆龍蝦免費射擊（DAY-199）
signal ice_phoenix(data: Dictionary)                   # 冰鳳凰覺醒 BOSS（DAY-200）
signal serial_bomb_crab(data: Dictionary)              # 連環炸彈蟹（DAY-201）
signal abyss_vortex(data: Dictionary)                  # 深淵漩渦魚（DAY-202）
signal humpback_whale(data: Dictionary)                # 座頭鯨覺醒（DAY-203）
signal free_spin_fish(data: Dictionary)                # 自由旋轉魚免費射擊（DAY-204）
signal jackpot_dragon(data: Dictionary)               # 獎池龍 Jackpot 抽獎（DAY-205）
signal comet_fish(data: Dictionary)                   # 彗星魚連鎖爆炸（DAY-206）
signal golden_wave_fish(data: Dictionary)             # 黃金波浪魚全場倍率衝擊（DAY-207）
signal dragon_king(data: Dictionary)                  # 深海龍王全服合力蓄力（DAY-208）
signal fortune_coin_fish(data: Dictionary)            # 幸運金幣魚即時獎勵（DAY-209）
signal lucky_hot_zone(data: Dictionary)               # 幸運熱區魚空間策略（DAY-210）
signal lucky_trident(data: Dictionary)               # 幸運三叉魚互動三轉盤（DAY-211）
signal time_freeze_fish(data: Dictionary)            # 時間凍結魚系統（DAY-212）
signal rainbow_prism(data: Dictionary)              # 彩虹稜鏡魚系統（DAY-213）
signal golden_accumulator(data: Dictionary)         # 黃金累積魚系統（DAY-214）
signal lucky_mirror_fish(data: Dictionary)          # 幸運鏡像魚系統（DAY-215）
signal cursed_poison_fish(data: Dictionary)         # 詛咒毒魚系統（DAY-216）
signal lucky_auction_fish(data: Dictionary)         # 幸運拍賣魚系統（DAY-217）
signal lucky_evolution_fish(data: Dictionary)       # 幸運進化魚系統（DAY-218）
signal lucky_infection_fish(data: Dictionary)       # 幸運連鎖感染魚系統（DAY-219）
signal lucky_ricochet_fish(data: Dictionary)        # 幸運反彈魚系統（DAY-220）
signal lucky_black_hole(data: Dictionary)           # 幸運黑洞魚系統（DAY-221）
signal lucky_resonance_fish(data: Dictionary)       # 幸運共鳴魚系統（DAY-222）
signal lucky_teleport_fish(data: Dictionary)        # 幸運傳送魚系統（DAY-223）
signal lucky_split_fish(data: Dictionary)           # 幸運分裂魚系統（DAY-224）
signal lucky_charge_fish(data: Dictionary)          # 幸運充能魚系統（DAY-225）
signal lucky_chain_bomb(data: Dictionary)           # 幸運鏈鎖爆炸魚系統（DAY-226）
signal lucky_mirror_time(data: Dictionary)          # 幸運鏡像時空魚系統（DAY-227）
signal lucky_quantum_fish(data: Dictionary)         # 幸運量子魚系統（DAY-228）
signal lucky_parasite_fish(data: Dictionary)        # 幸運寄生魚系統（DAY-229）
signal lucky_storm_fish(data: Dictionary)           # 幸運風暴魚系統（DAY-230）
signal lucky_boomerang_fish(data: Dictionary)       # 幸運迴旋鏢魚系統（DAY-231）
signal lucky_magnet_fish(data: Dictionary)          # 幸運磁力魚系統（DAY-232）
signal lucky_echo_fish(data: Dictionary)            # 幸運回聲魚系統（DAY-233）
signal lucky_vortex_fish(data: Dictionary)          # 幸運漩渦魚系統（DAY-234）
signal lucky_time_bomb_fish(data: Dictionary)       # 幸運時間炸彈魚系統（DAY-235）
signal lucky_mirror_world(data: Dictionary)         # 幸運鏡面世界魚系統（DAY-236）
signal lucky_freeze_world(data: Dictionary)         # 幸運冰凍世界魚系統（DAY-237）
signal lucky_gravity_flip(data: Dictionary)         # 幸運重力反轉魚系統（DAY-238）
signal royal_chain_lightning(chain_data: Dictionary)   # 皇家閃電鰻持續連鎖電擊（DAY-156）
signal golden_turtle_time_stop(data: Dictionary)       # 黃金海龜時間停止（DAY-159）
signal lucky_star_fish(data: Dictionary)               # 幸運星魚全場倍率翻倍（DAY-160）
signal golden_shark_berserk(data: Dictionary)          # 黃金鯊魚全服狂暴模式（DAY-161）
signal money_fish_reward(data: Dictionary)             # 金幣魚王即時獎勵（DAY-162）
signal captain_fish_race(data: Dictionary)             # 船長魚全服競速模式（DAY-163）
signal abyss_whale(data: Dictionary)                   # 深淵巨鯨全服 Boss 挑戰（DAY-164）
signal drill_lobster_chain(chain_data: Dictionary)      # 鑽頭龍蝦連帶效果（DAY-142）
signal bomb_crab_chain(chain_data: Dictionary)          # 炸彈蟹連環爆炸（DAY-143）
signal mega_octopus_wheel_start(wheel_data: Dictionary) # 巨型章魚轉盤開始（DAY-144）
signal mega_octopus_wheel_result(result_data: Dictionary) # 巨型章魚轉盤結果（DAY-144）
signal anglerfish_shock(shock_data: Dictionary)          # 鮟鱇魚電擊寶箱（DAY-145）
signal crocodile_hunt(hunt_data: Dictionary)             # 鱷魚獵魚累積（DAY-146）
signal giant_prize_fish(event_data: Dictionary)          # 夢幻巨型獎勵魚（DAY-147）
signal chainlong_wheel_start(wheel_data: Dictionary)     # 千龍王輪盤開始（DAY-148）
signal chainlong_wheel_result(result_data: Dictionary)   # 千龍王輪盤結果（DAY-148）
signal golden_jellyfish_shock(shock_data: Dictionary)    # 黃金水母全場電擊（DAY-149）
signal thunderbolt_lobster_activate(event_data: Dictionary) # 雷霆龍蝦免費射擊開始（DAY-150）
signal thunderbolt_lobster_shot(shot_data: Dictionary)      # 雷霆龍蝦自動射擊（DAY-150）
signal thunderbolt_lobster_end(end_data: Dictionary)        # 雷霆龍蝦免費射擊結束（DAY-150）
signal rainbow_phoenix_activate(event_data: Dictionary)     # 彩虹鳳凰 Power Up 開始（DAY-151）
signal rainbow_phoenix_end(end_data: Dictionary)            # 彩虹鳳凰 Power Up 結束（DAY-151）
signal vampire_grow(grow_data: Dictionary)                  # 吸血鬼倍率成長（DAY-152）
signal vampire_blood_moon(moon_data: Dictionary)            # 吸血鬼血月模式（DAY-152）
signal vampire_killed(kill_data: Dictionary)                # 吸血鬼被擊破（DAY-152）
signal crystal_dragon_drop(drop_data: Dictionary)           # 水晶龍掉落水晶（DAY-153）
signal crystal_dragon_reward(reward_data: Dictionary)       # 水晶龍地獄龍大獎（DAY-153）
signal crystal_dragon_status(status_data: Dictionary)       # 水晶龍狀態（DAY-153）
signal unlucky_bonus(bonus_data: Dictionary)           # 失敗補償觸發（DAY-135）
signal speed_race_started(race_data: Dictionary)       # 競速獵殺開始（DAY-136）
signal speed_race_ended(race_data: Dictionary)         # 競速獵殺結束（DAY-136）
signal speed_race_cancelled(race_data: Dictionary)     # 競速獵殺取消（DAY-136）
signal speed_race_result(result_data: Dictionary)      # 競速個人結果（DAY-136）
signal bounty_posted(bounty_data: Dictionary)          # 懸賞發布廣播（DAY-137）
signal bounty_claimed(claim_data: Dictionary)          # 懸賞個人領取通知（DAY-137）
signal bounty_killed(kill_data: Dictionary)            # 懸賞目標擊破廣播（DAY-137）
signal bounty_expired(expire_data: Dictionary)         # 懸賞過期通知（DAY-137）
signal mult_storm_started(storm_data: Dictionary)      # 倍率風暴開始（DAY-138）
signal mult_storm_ended(storm_data: Dictionary)        # 倍率風暴結束（DAY-138）
signal dual_roulette_started(roulette_data: Dictionary) # 雙環輪盤開始（DAY-139）
signal dual_roulette_result(result_data: Dictionary)    # 雙環輪盤結果（DAY-139）
signal mega_catch_started(event_data: Dictionary)       # Mega Catch 事件開始（DAY-140）
signal mega_catch_ended(event_data: Dictionary)         # Mega Catch 事件結束（DAY-140）
signal mystery_box_updated(box_data: Dictionary)       # 神秘寶箱狀態更新（DAY-090）
signal mystery_box_dropped(drop_data: Dictionary)      # 神秘寶箱掉落通知（DAY-090）
signal mystery_box_opened(open_data: Dictionary)       # 神秘寶箱開箱結果（DAY-090）
# 房間難度系統（DAY-091）
signal room_list_received(room_data: Dictionary)       # 房間列表更新
signal room_switched(switch_data: Dictionary)          # 房間切換成功
signal room_error(error_data: Dictionary)              # 房間操作失敗
# 每日簽到轉盤（DAY-092）
signal daily_spin_state(state_data: Dictionary)        # 每日轉盤狀態
signal daily_spin_result(result_data: Dictionary)      # 每日轉盤結果
# 商店系統（DAY-094）
signal shop_updated(shop_data: Dictionary)             # 商店狀態更新
signal shop_purchased(purchase_data: Dictionary)       # 購買成功通知
signal shop_error(error_data: Dictionary)              # 購買失敗通知
# 名人堂系統（DAY-110）
signal hall_of_fame_updated(hall_data: Dictionary)     # 名人堂更新
signal hall_of_fame_new_record(record_data: Dictionary) # 新記錄誕生廣播
# 智慧推薦系統（DAY-110）
signal recommendations_received(rec_data: Dictionary)  # 推薦結果

# 遊戲狀態
var current_state: String = "normal_play"
var player_data: Dictionary = {}
var targets: Dictionary = {}  # instance_id -> target_data

# 角色顏色對應（規格書 5章）
const CHARACTER_COLORS = {
	"chiikawa": Color(1.0, 0.6, 0.8),   # 粉紅
	"hachiware": Color(0.4, 0.6, 1.0),  # 藍色
	"usagi": Color(1.0, 0.9, 0.2),      # 黃色
}

# 角色名稱
const CHARACTER_NAMES = {
	"chiikawa": "Chiikawa",
	"hachiware": "Hachiware",
	"usagi": "Usagi",
}

func _ready() -> void:
	# 連接 NetworkManager 訊號
	NetworkManager.message_received.connect(_on_message_received)
	NetworkManager.connected.connect(_on_connected)
	NetworkManager.disconnected.connect(_on_disconnected)
	# 啟動資產預載入（背景執行，不阻塞遊戲）
	call_deferred("_start_preloading")

func _start_preloading() -> void:
	if LoadingManager != null:
		LoadingManager.preload_all()

func _on_connected() -> void:
	print("[GameManager] Connected to server")

func _on_disconnected() -> void:
	print("[GameManager] Disconnected from server")

## 處理 Server 訊息
func _on_message_received(type: String, payload: Dictionary) -> void:
	match type:
		"game_state":
			_handle_game_state(payload)
		"target_spawn":
			_handle_target_spawn(payload)
		"target_update":
			_handle_target_update(payload)
		"target_kill":
			_handle_target_kill(payload)
		"attack_result":
			_handle_attack_result(payload)
		"reward":
			_handle_reward(payload)
		"player_update":
			_handle_player_update(payload)
		"boss_event":
			_handle_boss_event(payload)
		"bonus_event":
			_handle_bonus_event(payload)
		"leaderboard":
			_handle_leaderboard(payload)
		"achievement":
			_handle_achievement(payload)
		"combo_event":
			_handle_combo_event(payload)
		"mission_update":
			_handle_mission_update(payload)
		"mission_complete":
			_handle_mission_complete(payload)
		"jackpot_update":
			_handle_jackpot_update(payload)
		"jackpot_win":
			_handle_jackpot_win(payload)
		"jackpot_animation":
			_handle_jackpot_animation(payload)
		"player_stats_update":
			_handle_player_stats_update(payload)
		"announcement":
			_handle_announcement(payload)
		"daily_bonus":
			_handle_daily_bonus(payload)
		"spectator_join":
			_handle_spectator_join(payload)
		"spectator_leave":
			_handle_spectator_leave(payload)
		"tournament_update":
			_handle_tournament_update(payload)
		# 每日賽排名更新（DAY-093）
		"daily_tournament_update":
			_handle_daily_tournament_update(payload)
		# 多格式賽排名更新（DAY-111）
		"multi_format_update":
			_handle_multi_format_update(payload)
		# 動態牆事件（DAY-112）
		"activity_feed_event":
			emit_signal("activity_feed_event", payload)
		"activity_feed_history":
			emit_signal("activity_feed_history", payload)
		# 雙層倍率輪盤（DAY-113）
		"roulette_start":
			_handle_roulette_start(payload)
		"roulette_result":
			_handle_roulette_result(payload)
		# Buy Bonus 系統（DAY-114）
		"buy_bonus_success":
			emit_signal("buy_bonus_success", payload)
		"buy_bonus_error":
			emit_signal("buy_bonus_error", payload)
		"buy_bonus_status":
			emit_signal("buy_bonus_status", payload)
		# Co-op Boss Raid 系統（DAY-115）
		"raid_warning":
			emit_signal("raid_warning", payload)
		"raid_start":
			emit_signal("raid_started", payload)
		"raid_update":
			emit_signal("raid_updated", payload)
		"raid_result":
			emit_signal("raid_result", payload)
		# 碎片收集大獎系統（DAY-116）
		"fragment_drop":
			emit_signal("fragment_dropped", payload)
		"fragment_complete":
			# 廣播時 Client 端判斷是否為自己
			var my_id = get_player_id() if has_method("get_player_id") else ""
			payload["is_self"] = (payload.get("player_id", "") == my_id)
			emit_signal("fragment_completed", payload)
		"fragment_status":
			emit_signal("fragment_status_received", payload)
		"title_unlocked":
			_handle_title_unlocked(payload)
		"skin_update":
			_handle_skin_update(payload)
		"season_update":
			_handle_season_update(payload)
		"season_level_up":
			_handle_season_level_up(payload)
		"friend_list":
			_handle_friend_list(payload)
		"friend_request":
			_handle_friend_request(payload)
		"friend_update":
			_handle_friend_update(payload)
		# 好友禮物系統（DAY-101）
		"gift_received":
			_handle_gift_received(payload)
		"gift_sent":
			_handle_gift_sent(payload)
		"gift_status":
			_handle_gift_status(payload)
		"gift_error":
			_handle_gift_error(payload)
		# 好友挑戰系統（DAY-102）
		"challenge_request":
			_handle_challenge_request(payload)
		"challenge_update":
			_handle_challenge_updated(payload)
		"challenge_result":
			_handle_challenge_result(payload)
		"challenge_error":
			_handle_challenge_error(payload)
		# 私訊系統（DAY-103）
		"dm_received":
			_handle_dm_received(payload)
		"dm_sent":
			_handle_dm_sent(payload)
		"dm_error":
			_handle_dm_error(payload)
		# 玩家名片系統（DAY-106）
		"player_card":
			emit_signal("player_card_received", payload)
		# 登入里程碑系統（DAY-107）
		"login_milestone":
			_handle_login_milestone(payload)
		"login_progress":
			_handle_login_progress(payload)
		# 超級 Bonus 系統（DAY-108）
		"super_bonus_ready":
			_handle_super_bonus_ready(payload)
		"guild_update":
			_handle_guild_update(payload)
		"guild_list":
			_handle_guild_list(payload)
		"guild_task_complete":
			_handle_guild_task_complete(payload)
		"guild_error":
			_handle_guild_error(payload)
		"guild_message":
			_handle_guild_message(payload)
		"guild_war_update":
			_handle_guild_war_update(payload)
		"guild_war_result":
			_handle_guild_war_result(payload)
		"daily_boss_update":
			_handle_daily_boss_update(payload)
		"daily_boss_defeated":
			_handle_daily_boss_defeated(payload)
		"vip_update":
			_handle_vip_update(payload)
		"vip_level_up":
			_handle_vip_level_up(payload)
		"vip_weekly_claimed":
			_handle_vip_weekly_claimed(payload)
		"event_update":
			_handle_event_update(payload)
		"codex_update":
			_handle_codex_update(payload)
		"codex_unlock":
			_handle_codex_unlock(payload)
		"codex_complete":
			_handle_codex_complete(payload)
		"referral_info":
			_handle_referral_info(payload)
		"referral_success":
			_handle_referral_success(payload)
		"referral_error":
			_handle_referral_error(payload)
		"wheel_trigger":
			_handle_wheel_trigger(payload)
		"challenge_unlocked":
			_handle_challenge_unlocked(payload)
		"mission_streak_bonus":
			_handle_mission_streak_bonus(payload)
		"weather_update":
			_handle_weather_update(payload)
		"chain_explosion":
			_handle_chain_explosion(payload)
		"special_weapon_update":
			_handle_special_weapon_update(payload)
		"special_weapon_fired":
			_handle_special_weapon_fired(payload)
		"special_weapon_charged":
			_handle_special_weapon_charged(payload)
		"homing_missile_result":
			_handle_homing_missile_result(payload)
		"dragon_wrath_charge":
			_handle_dragon_wrath_charge(payload)
		"dragon_wrath_result":
			_handle_dragon_wrath_result(payload)
		"torpedo_result":
			_handle_torpedo_result(payload)
		"railgun_result":
			_handle_railgun_result(payload)
		"black_hole_result":
			_handle_black_hole_result(payload)
		"roulette_crab_start":
			_handle_roulette_crab_start(payload)
		"roulette_crab_result":
			_handle_roulette_crab_result(payload)
		"lion_dance_burst":
			_handle_lion_dance_burst(payload)
		"vortex_fish":
			_handle_vortex_fish(payload)
		"freeze_bomb":
			_handle_freeze_bomb(payload)
		"ice_fishing_wheel":
			_handle_ice_fishing_wheel(payload)
		"lucky_egg_fish":
			_handle_lucky_egg_fish(payload)
		"rainbow_lucky_fish":
			_handle_rainbow_lucky_fish(payload)
		"sea_anemone":
			_handle_sea_anemone(payload)
		"lucky_dice_fish":
			_handle_lucky_dice_fish(payload)
		"fire_storm_fish":
			_handle_fire_storm_fish(payload)
		"golden_treasure_fish":
			_handle_golden_treasure_fish(payload)
		"mermaid_healing":
			_handle_mermaid_healing(payload)
		"lucky_clover_fish":
			_handle_lucky_clover_fish(payload)
		"rainbow_shark_burst":
			_handle_rainbow_shark_burst(payload)
		"thunder_shark_chain":
			_handle_thunder_shark_chain(payload)
		"vampire_fish":
			_handle_vampire_fish(payload)
		"lightning_auto_chain":
			_handle_lightning_auto_chain(payload)
		"meteor_fish":
			_handle_meteor_fish(payload)
		"phoenix_fish":
			_handle_phoenix_fish(payload)
		"dragon_turtle":
			_handle_dragon_turtle(payload)
		"chain_bomb":
			_handle_chain_bomb(payload)
		"crocodile_hunter":
			_handle_crocodile_hunter(payload)
		"time_bomb_fish":
			_handle_time_bomb_fish(payload)
		"triple_lucky_fish":
			_handle_triple_lucky_fish(payload)
		"school_panic":
			_handle_school_panic(payload)
		"rock_skeleton_concert":
			_handle_rock_skeleton_concert(payload)
		"electric_jellyfish":
			_handle_electric_jellyfish(payload)
		"chainlong_king":
			_handle_chainlong_king(payload)
		"drill_lobster":
			_handle_drill_lobster(payload)
		"anglerfish_electric":
			_handle_anglerfish_electric(payload)
		"mystic_dragon":
			_handle_mystic_dragon(payload)
		"ghost_fish":
			_handle_ghost_fish(payload)
		"thunderbolt_lobster":
			_handle_thunderbolt_lobster(payload)
		"ice_phoenix":
			_handle_ice_phoenix(payload)
		"serial_bomb_crab":
			_handle_serial_bomb_crab(payload)
		"abyss_vortex":
			_handle_abyss_vortex(payload)
		"humpback_whale":
			_handle_humpback_whale(payload)
		"free_spin_fish":
			_handle_free_spin_fish(payload)
		"jackpot_dragon":
			_handle_jackpot_dragon(payload)
		"comet_fish":
			_handle_comet_fish(payload)
		"golden_wave_fish":
			_handle_golden_wave_fish(payload)
		"dragon_king":
			_handle_dragon_king(payload)
		"fortune_coin_fish":
			_handle_fortune_coin_fish(payload)
		"lucky_hot_zone":
			_handle_lucky_hot_zone(payload)
		"lucky_trident":
			_handle_lucky_trident(payload)
		"time_freeze_fish":
			_handle_time_freeze_fish(payload)
		"rainbow_prism":
			_handle_rainbow_prism(payload)
		"golden_accumulator":
			_handle_golden_accumulator(payload)
		"lucky_mirror_fish":
			_handle_lucky_mirror_fish(payload)
		"cursed_poison_fish":
			_handle_cursed_poison_fish(payload)
		"lucky_auction_fish":
			_handle_lucky_auction_fish(payload)
		"lucky_evolution_fish":
			_handle_lucky_evolution_fish(payload)
		"lucky_infection_fish":
			_handle_lucky_infection_fish(payload)
		"lucky_ricochet_fish":
			_handle_lucky_ricochet_fish(payload)
		"lucky_black_hole":
			_handle_lucky_black_hole(payload)
		"lucky_resonance_fish":
			_handle_lucky_resonance_fish(payload)
		"lucky_teleport_fish":
			_handle_lucky_teleport_fish(payload)
		"lucky_split_fish":
			_handle_lucky_split_fish(payload)
		"lucky_charge_fish":
			_handle_lucky_charge_fish(payload)
		"lucky_chain_bomb":
			_handle_lucky_chain_bomb(payload)
		"lucky_mirror_time":
			_handle_lucky_mirror_time(payload)
		"lucky_quantum_fish":
			_handle_lucky_quantum_fish(payload)
		"lucky_parasite_fish":
			_handle_lucky_parasite_fish(payload)
		"lucky_storm_fish":
			_handle_lucky_storm_fish(payload)
		"lucky_boomerang_fish":
			_handle_lucky_boomerang_fish(payload)
		"lucky_magnet_fish":
			_handle_lucky_magnet_fish(payload)
		"lucky_echo_fish":
			_handle_lucky_echo_fish(payload)
		"lucky_vortex_fish":
			_handle_lucky_vortex_fish(payload)
		"lucky_time_bomb_fish":
			_handle_lucky_time_bomb_fish(payload)
		"lucky_mirror_world":
			_handle_lucky_mirror_world(payload)
		"lucky_freeze_world":
			_handle_lucky_freeze_world(payload)
		"lucky_gravity_flip":
			_handle_lucky_gravity_flip(payload)
		"golden_turtle_time_stop":
			_handle_golden_turtle_time_stop(payload)
		"lucky_star_fish":
			_handle_lucky_star_fish(payload)
		"golden_shark_berserk":
			_handle_golden_shark_berserk(payload)
		"money_fish_reward":
			_handle_money_fish_reward(payload)
		"captain_fish_race":
			_handle_captain_fish_race(payload)
		"abyss_whale":
			_handle_abyss_whale(payload)
		"royal_chain_lightning":
			_handle_royal_chain_lightning(payload)
		"drill_lobster_chain":
			_handle_drill_lobster_chain(payload)
		"bomb_crab_chain":
			_handle_bomb_crab_chain(payload)
		"mega_octopus_wheel_start":
			emit_signal("mega_octopus_wheel_start", payload)
		"mega_octopus_wheel_result":
			emit_signal("mega_octopus_wheel_result", payload)
		"anglerfish_shock":
			emit_signal("anglerfish_shock", payload)
		"crocodile_hunt":
			emit_signal("crocodile_hunt", payload)
		"giant_prize_fish":
			emit_signal("giant_prize_fish", payload)
		# 千龍王強化輪盤系統（DAY-148）
		"chainlong_wheel_start":
			_handle_chainlong_wheel_start(payload)
		"chainlong_wheel_result":
			_handle_chainlong_wheel_result(payload)
		"chainlong_wheel_status":
			pass  # 冷卻狀態，目前不需要特別處理
		# 黃金水母全場電擊系統（DAY-149）
		"golden_jellyfish_shock":
			emit_signal("golden_jellyfish_shock", payload)
		# 雷霆龍蝦免費射擊系統（DAY-150）
		"thunderbolt_lobster_activate":
			emit_signal("thunderbolt_lobster_activate", payload)
		"thunderbolt_lobster_shot":
			emit_signal("thunderbolt_lobster_shot", payload)
		"thunderbolt_lobster_end":
			emit_signal("thunderbolt_lobster_end", payload)
		# 彩虹鳳凰 Power Up 系統（DAY-151）
		"rainbow_phoenix_activate":
			emit_signal("rainbow_phoenix_activate", payload)
		"rainbow_phoenix_end":
			emit_signal("rainbow_phoenix_end", payload)
		"rainbow_phoenix_status":
			pass  # 狀態，目前不需要特別處理
		# 吸血鬼成長倍率系統（DAY-152）
		"vampire_grow":
			emit_signal("vampire_grow", payload)
		"vampire_blood_moon":
			emit_signal("vampire_blood_moon", payload)
		"vampire_killed":
			emit_signal("vampire_killed", payload)
		# 水晶龍收集大獎系統（DAY-153）
		"crystal_dragon_drop":
			emit_signal("crystal_dragon_drop", payload)
		"crystal_dragon_update":
			emit_signal("crystal_dragon_status", payload)
		"crystal_dragon_reward":
			emit_signal("crystal_dragon_reward", payload)
		"crystal_dragon_status":
			emit_signal("crystal_dragon_status", payload)
		"unlucky_bonus":
			_handle_unlucky_bonus(payload)
		"speed_race_start":
			_handle_speed_race_start(payload)
		"speed_race_end":
			_handle_speed_race_end(payload)
		"speed_race_cancel":
			_handle_speed_race_cancel(payload)
		"speed_race_result":
			_handle_speed_race_result(payload)
		"bounty_posted":
			_handle_bounty_posted(payload)
		"bounty_claimed":
			_handle_bounty_claimed(payload)
		"bounty_killed":
			_handle_bounty_killed(payload)
		"bounty_expired":
			_handle_bounty_expired(payload)
		"mult_storm_start":
			_handle_mult_storm_start(payload)
		"mult_storm_end":
			_handle_mult_storm_end(payload)
		"dual_roulette_start":
			_handle_dual_roulette_start(payload)
		"dual_roulette_result":
			_handle_dual_roulette_result(payload)
		"dual_roulette_status":
			pass  # 冷卻狀態，目前不需要特別處理
		"mega_catch_start":
			_handle_mega_catch_start(payload)
		"mega_catch_end":
			_handle_mega_catch_end(payload)
		"mega_catch_status":
			pass  # 登入時狀態，由 MegaCatchPanel 處理
		"mystery_box_drop":
			_handle_mystery_box_drop(payload)
		"mystery_box_update":
			_handle_mystery_box_update(payload)
		"mystery_box_opened":
			_handle_mystery_box_opened(payload)
		"streak_update":
			_handle_streak_update(payload)
		"streak_reset":
			_handle_streak_reset(payload)
		# 房間難度系統（DAY-091）
		"room_list":
			_handle_room_list(payload)
		"room_switched":
			_handle_room_switched(payload)
		"room_error":
			_handle_room_error(payload)
		# 每日簽到轉盤（DAY-092）
		"daily_spin_state":
			_handle_daily_spin_state(payload)
		"daily_spin_result":
			_handle_daily_spin_result(payload)
		# 商店系統（DAY-094）
		"shop_update":
			_handle_shop_update(payload)
		"shop_purchased":
			_handle_shop_purchased(payload)
		"shop_error":
			_handle_shop_error(payload)
		"error":
			_handle_error(payload)
		"pong":
			pass  # Ping/Pong 心跳
		# 賽季節日活動系統（DAY-109）
		"festival_update":
			_handle_festival_update(payload)
		"festival_task_ready":
			_handle_festival_task_ready(payload)
		"festival_task_claimed":
			_handle_festival_task_claimed(payload)
		"festival_title_earned":
			_handle_festival_title_earned(payload)
		"festival_error":
			_handle_festival_error(payload)
		# 名人堂系統（DAY-110）
		"hall_of_fame_update":
			_handle_hall_of_fame_update(payload)
		"hall_of_fame_new_record":
			_handle_hall_of_fame_new_record(payload)
		# 智慧推薦系統（DAY-110）
		"recommendations":
			_handle_recommendations(payload)
		# 幸運捕獲系統（DAY-119）
		"lucky_catch":
			emit_signal("lucky_catch", payload)
		# 任務連續寬限期（DAY-120）
		"mission_mercy_protected":
			emit_signal("mission_mercy_protected", payload)
		# Rapid Respin 系統（DAY-121）
		"rapid_respin":
			emit_signal("rapid_respin", payload)
		"rapid_respin_end":
			emit_signal("rapid_respin_end", payload)
		# 寶藏地圖系統（DAY-122）
		"treasure_map_update":
			emit_signal("treasure_map_updated", payload)
		"treasure_map_line":
			emit_signal("treasure_map_line", payload)
		"treasure_map_full":
			emit_signal("treasure_map_full", payload)
		# 閃電挑戰系統（DAY-123）
		"flash_challenge_start":
			emit_signal("flash_challenge_started", payload)
		"flash_challenge_update":
			emit_signal("flash_challenge_updated", payload)
		"flash_challenge_end":
			emit_signal("flash_challenge_ended", payload)
		"flash_challenge_reward":
			emit_signal("flash_challenge_reward", payload)
		# 傳說目標警報系統（DAY-124）
		"rare_target_alert":
			emit_signal("rare_target_alerted", payload)
		# 黃金時間系統（DAY-125）
		"golden_time_start":
			emit_signal("golden_time_started", payload)
		"golden_time_end":
			emit_signal("golden_time_ended", payload)
		"golden_time_status":
			emit_signal("golden_time_status", payload)
		# 稀有連擊累積倍率系統（DAY-126）
		"rare_catch_update":
			emit_signal("rare_catch_updated", payload)
		"rare_catch_broadcast":
			emit_signal("rare_catch_broadcasted", payload)
		"rare_catch_reset":
			emit_signal("rare_catch_reset", payload)
		# 天氣湧現事件（DAY-127）
		"weather_surge_start":
			emit_signal("weather_surge_started", payload)
		"weather_surge_end":
			emit_signal("weather_surge_ended", payload)
		# 龍怒蓄力大招系統（DAY-128）
		"wrath_update":
			emit_signal("wrath_updated", payload)
		"wrath_start":
			emit_signal("wrath_started", payload)
		"wrath_result":
			emit_signal("wrath_result", payload)
		# 不死 BOSS 連勝系統（DAY-129）
		"immortal_boss_spawn":
			emit_signal("immortal_boss_spawned", payload)
		"immortal_boss_hit":
			emit_signal("immortal_boss_hit", payload)
		"immortal_boss_leave":
			emit_signal("immortal_boss_left", payload)
		"immortal_boss_status":
			emit_signal("immortal_boss_status", payload)
		# 覺醒 BOSS 系統（DAY-130）
		"awaken_boss_spawn":
			emit_signal("awaken_boss_spawned", payload)
		"awaken_boss_hit":
			emit_signal("awaken_boss_hit", payload)
		"awaken_boss_powerup":
			emit_signal("awaken_boss_powerup", payload)
		"awaken_boss_leave":
			emit_signal("awaken_boss_left", payload)
		"awaken_boss_status":
			emit_signal("awaken_boss_status", payload)
		# 連勝獎勵系統（DAY-131）
		"win_streak_update":
			emit_signal("win_streak_updated", payload)
		"win_streak_milestone":
			emit_signal("win_streak_milestone", payload)
		"win_streak_reset":
			emit_signal("win_streak_reset", payload)
		# 閃電鰻連鎖攻擊系統（DAY-132）
		"lightning_eel_chain":
			emit_signal("lightning_eel_chain", payload)
		"lightning_eel_status":
			emit_signal("lightning_eel_status", payload)
		# 狂熱模式系統（DAY-133）
		"fever_mode_start":
			emit_signal("fever_mode_started", payload)
		"fever_mode_end":
			emit_signal("fever_mode_ended", payload)
		"fever_mode_status":
			emit_signal("fever_mode_status", payload)

func _handle_game_state(payload: Dictionary) -> void:
	var new_state = payload.get("state", "")
	if new_state != current_state:
		current_state = new_state
		print("[GameManager] State: ", current_state)
		emit_signal("game_state_changed", current_state)

func _handle_target_spawn(payload: Dictionary) -> void:
	var instance_id = payload.get("instance_id", "")
	targets[instance_id] = payload
	emit_signal("target_spawned", payload)

func _handle_target_update(payload: Dictionary) -> void:
	var instance_id = payload.get("instance_id", "")
	if targets.has(instance_id):
		targets[instance_id].merge(payload, true)
	emit_signal("target_updated", payload)

func _handle_target_kill(payload: Dictionary) -> void:
	var instance_id = payload.get("instance_id", "")
	targets.erase(instance_id)
	emit_signal("target_killed", payload)

func _handle_attack_result(payload: Dictionary) -> void:
	emit_signal("attack_result", payload)

func _handle_reward(payload: Dictionary) -> void:
	# 更新本地金幣顯示
	if payload.has("new_balance"):
		player_data["coins"] = payload["new_balance"]
	emit_signal("reward_received", payload)

func _handle_player_update(payload: Dictionary) -> void:
	player_data = payload
	emit_signal("player_updated", payload)

func _handle_boss_event(payload: Dictionary) -> void:
	emit_signal("boss_event", payload)
	# BGM 切換：BOSS Phase 2 時切換到 boss_rage
	var event = payload.get("event", "")
	if event == "phase_change" and payload.get("phase", 1) == 2:
		if AudioManager != null:
			AudioManager.play_bgm(AudioManager.BGM.BOSS_RAGE)
	elif event == "kill":
		# BOSS 擊敗：短暫靜音，等 boss_result 狀態切換回主 BGM
		if AudioManager != null:
			AudioManager.stop_bgm_briefly()

func _handle_bonus_event(payload: Dictionary) -> void:
	emit_signal("bonus_event", payload)

func _handle_leaderboard(payload: Dictionary) -> void:
	var entries = payload.get("entries", [])
	emit_signal("leaderboard_updated", entries)

func _handle_achievement(payload: Dictionary) -> void:
	print("[GameManager] Achievement unlocked: ", payload.get("name", ""))
	emit_signal("achievement_unlocked", payload)

func _handle_combo_event(payload: Dictionary) -> void:
	## 連擊事件（DAY-022）
	var combo_count = payload.get("combo_count", 1)
	print("[GameManager] COMBO x%d!" % combo_count)
	emit_signal("combo_event", payload)

## 任務進度更新（DAY-037）
func _handle_mission_update(payload: Dictionary) -> void:
	var missions = payload.get("missions", [])
	var reset_at = payload.get("reset_at", 0)
	emit_signal("mission_updated", missions)
	# 傳遞重置時間給 HUD（DAY-038）
	if reset_at > 0:
		var hud = get_node_or_null("/root/Main/HUD")
		if is_instance_valid(hud) and hud.has_method("set_mission_reset_at"):
			hud.set_mission_reset_at(reset_at)

## 任務完成通知（DAY-037）
func _handle_mission_complete(payload: Dictionary) -> void:
	print("[GameManager] Mission completed: ", payload.get("name", ""))
	emit_signal("mission_completed", payload)

func _handle_error(payload: Dictionary) -> void:
	push_warning("[GameManager] Server error: " + str(payload))

## Jackpot 池更新（DAY-048）
func _handle_jackpot_update(payload: Dictionary) -> void:
	emit_signal("jackpot_updated", payload)

## Jackpot 中獎通知（DAY-048）
func _handle_jackpot_win(payload: Dictionary) -> void:
	var level = payload.get("level", "mini")
	var amount = payload.get("amount", 0)
	var winner_name = payload.get("winner_name", "")
	print("[GameManager] JACKPOT WIN! Level=%s Amount=%d Winner=%s" % [level, amount, winner_name])
	emit_signal("jackpot_won", payload)

## Jackpot 觸發動畫通知（DAY-095）
func _handle_jackpot_animation(payload: Dictionary) -> void:
	var level = payload.get("level", "mini")
	var amount = payload.get("amount", 0)
	print("[GameManager] JACKPOT ANIMATION! Level=%s Amount=%d" % [level, amount])
	emit_signal("jackpot_animation", payload)

## 玩家統計更新（DAY-096）
func _handle_player_stats_update(payload: Dictionary) -> void:
	emit_signal("player_stats_updated", payload)

## 全服公告（DAY-097）
func _handle_announcement(payload: Dictionary) -> void:
	emit_signal("announcement_received", payload)

## 觀戰者加入通知（DAY-054d）
func _handle_spectator_join(payload: Dictionary) -> void:
	var count = payload.get("spectator_count", 1)
	print("[GameManager] Spectator joined! Total spectators: %d" % count)
	emit_signal("spectator_joined", payload)

## 每日登入獎勵（DAY-065）
func _handle_daily_bonus(payload: Dictionary) -> void:
	var streak = payload.get("streak", 1)
	var reward = payload.get("reward", 0)
	var is_new = payload.get("is_new_streak", false)
	if is_new:
		print("[GameManager] Daily bonus! Streak=%d Reward=%d" % [streak, reward])
		emit_signal("daily_bonus_received", payload)

## 觀戰者離開通知（DAY-055）
func _handle_spectator_leave(payload: Dictionary) -> void:
	var count = payload.get("spectator_count", 0)
	print("[GameManager] Spectator left! Remaining spectators: %d" % count)
	emit_signal("spectator_left", payload)

## 週賽排名更新（DAY-066）
func _handle_tournament_update(payload: Dictionary) -> void:
	var rank = payload.get("player_rank", 0)
	var points = payload.get("player_points", 0)
	if rank > 0:
		print("[GameManager] Tournament rank=%d points=%d" % [rank, points])
	emit_signal("tournament_updated", payload)

## 每日賽排名更新（DAY-093）
func _handle_daily_tournament_update(payload: Dictionary) -> void:
	var rank = payload.get("player_rank", 0)
	var points = payload.get("player_points", 0)
	if rank > 0:
		print("[GameManager] Daily Tournament rank=%d points=%d" % [rank, points])
	emit_signal("daily_tournament_updated", payload)

## 多格式賽排名更新（DAY-111）
func _handle_multi_format_update(payload: Dictionary) -> void:
	var rank = payload.get("player_rank", 0)
	var score = payload.get("player_score", 0.0)
	var format_name = payload.get("format_name", "")
	if rank > 0:
		print("[GameManager] MultiFormat rank=%d score=%.1f format=%s" % [rank, score, format_name])
	emit_signal("multi_format_updated", payload)

## 商店狀態更新（DAY-094）
func _handle_shop_update(payload: Dictionary) -> void:
	emit_signal("shop_updated", payload)

## 購買成功通知（DAY-094）
func _handle_shop_purchased(payload: Dictionary) -> void:
	var item_name = payload.get("item_name", "")
	print("[GameManager] Shop purchased: %s" % item_name)
	emit_signal("shop_purchased", payload)

## 購買失敗通知（DAY-094）
func _handle_shop_error(payload: Dictionary) -> void:
	var reason = payload.get("reason", "")
	print("[GameManager] Shop error: %s" % reason)
	emit_signal("shop_error", payload)

## 處理稱號解鎖通知（DAY-068）
func _handle_title_unlocked(payload: Dictionary) -> void:
	emit_signal("title_unlocked", payload)
	# 播放稱號解鎖音效（用 big_win 音效）
	if AudioManager != null:
		AudioManager.play_sfx("big_win")

## 處理砲台外觀更新（DAY-071）
func _handle_skin_update(payload: Dictionary) -> void:
	emit_signal("skin_updated", payload)

## 處理賽季通行證更新（DAY-072）
func _handle_season_update(payload: Dictionary) -> void:
	emit_signal("season_updated", payload)

## 處理賽季等級升級（DAY-072）
func _handle_season_level_up(payload: Dictionary) -> void:
	emit_signal("season_level_up", payload)

## 處理好友列表更新（DAY-073）
func _handle_friend_list(payload: Dictionary) -> void:
	emit_signal("friend_list_updated", payload)

## 處理好友請求通知（DAY-073）
func _handle_friend_request(payload: Dictionary) -> void:
	emit_signal("friend_request_received", payload)

## 處理好友狀態更新（DAY-073）
func _handle_friend_update(payload: Dictionary) -> void:
	emit_signal("friend_updated", payload)

## 處理禮物收到通知（DAY-101）
func _handle_gift_received(payload: Dictionary) -> void:
	emit_signal("gift_received", payload)

## 處理禮物送出成功通知（DAY-101）
func _handle_gift_sent(payload: Dictionary) -> void:
	emit_signal("gift_sent", payload)

## 處理禮物狀態（DAY-101）
func _handle_gift_status(payload: Dictionary) -> void:
	emit_signal("gift_status", payload)

## 處理禮物錯誤（DAY-101）
func _handle_gift_error(payload: Dictionary) -> void:
	emit_signal("gift_error", payload)

## 處理挑戰邀請（DAY-102）
func _handle_challenge_request(payload: Dictionary) -> void:
	emit_signal("challenge_request", payload)

## 處理挑戰狀態更新（DAY-102）
func _handle_challenge_updated(payload: Dictionary) -> void:
	emit_signal("challenge_updated", payload)

## 處理挑戰結果（DAY-102）
func _handle_challenge_result(payload: Dictionary) -> void:
	emit_signal("challenge_result", payload)

## 處理挑戰錯誤（DAY-102）
func _handle_challenge_error(payload: Dictionary) -> void:
	emit_signal("challenge_error", payload)

## 處理收到私訊（DAY-103）
func _handle_dm_received(payload: Dictionary) -> void:
	emit_signal("dm_received", payload)

## 處理私訊發送確認（DAY-103）
func _handle_dm_sent(payload: Dictionary) -> void:
	emit_signal("dm_sent", payload)

## 處理私訊發送失敗（DAY-103）
func _handle_dm_error(payload: Dictionary) -> void:
	emit_signal("dm_error", payload)

## 處理公會資訊更新（DAY-074）
func _handle_guild_update(payload: Dictionary) -> void:
	emit_signal("guild_updated", payload)

## 處理公會列表（DAY-074）
func _handle_guild_list(payload: Dictionary) -> void:
	# 如果沒有公會，自動加入第一個公會（簡化版）
	var guilds: Array = payload.get("guilds", [])
	if guilds.size() > 0:
		var first_guild: Dictionary = guilds[0]
		var guild_id: String = first_guild.get("guild_id", "")
		if guild_id != "":
			send_message("join_guild", {"guild_id": guild_id})

## 處理公會任務完成（DAY-074）
func _handle_guild_task_complete(payload: Dictionary) -> void:
	emit_signal("guild_task_complete", payload)

## 處理公會錯誤（DAY-074）
func _handle_guild_error(payload: Dictionary) -> void:
	var msg: String = payload.get("message", "公會操作失敗")
	print("[GameManager] Guild error: ", msg)

## 處理公會聊天訊息（DAY-075）
func _handle_guild_message(payload: Dictionary) -> void:
	emit_signal("guild_message_received", payload)

## 處理公會戰排名更新（DAY-076）
func _handle_guild_war_update(payload: Dictionary) -> void:
	emit_signal("guild_war_updated", payload)

## 處理公會戰結算通知（DAY-076）
func _handle_guild_war_result(payload: Dictionary) -> void:
	emit_signal("guild_war_result", payload)

## 請求公會戰狀態（DAY-076）
func request_guild_war_status() -> void:
	NetworkManager.send_message("get_guild_war_status", {})

## 處理每日 BOSS 狀態更新（DAY-077）
func _handle_daily_boss_update(payload: Dictionary) -> void:
	emit_signal("daily_boss_updated", payload)

## 處理每日 BOSS 擊殺通知（DAY-077）
func _handle_daily_boss_defeated(payload: Dictionary) -> void:
	emit_signal("daily_boss_defeated", payload)

## 請求每日 BOSS 狀態（DAY-077）
func request_daily_boss_status() -> void:
	NetworkManager.send_message("get_daily_boss", {})

## 取得顯示名稱
func get_display_name() -> String:
	return player_data.get("display_name", "玩家")

## 取得目前角色顏色
func get_character_color() -> Color:
	var char_id = player_data.get("character_id", "chiikawa")
	return CHARACTER_COLORS.get(char_id, Color.WHITE)

## 取得目前角色名稱
func get_character_name() -> String:
	var char_id = player_data.get("character_id", "chiikawa")
	return CHARACTER_NAMES.get(char_id, "吉伊卡哇")

## 取得目前金幣
func get_coins() -> int:
	return player_data.get("coins", 0)

## 取得勞動值
func get_labor_value() -> int:
	return player_data.get("labor_value", 0)

## 取得投注等級
func get_bet_level() -> int:
	return player_data.get("bet_level", 1)

## 取得投注消耗
func get_bet_cost() -> int:
	return player_data.get("bet_cost", 1)

## 是否自動攻擊
func is_auto() -> bool:
	return player_data.get("is_auto", false)

## 取得鎖定目標 ID
func get_lock_target_id() -> String:
	return player_data.get("lock_target_id", "")

## 取得玩家 ID（用於排行榜標記自己）
func get_player_id() -> String:
	return player_data.get("id", "")

## 處理 VIP 狀態更新（DAY-078）
func _handle_vip_update(payload: Dictionary) -> void:
	emit_signal("vip_updated", payload)

## 處理 VIP 升級通知（DAY-078）
func _handle_vip_level_up(payload: Dictionary) -> void:
	emit_signal("vip_level_up", payload)

## 處理 VIP 週獎勵領取通知（DAY-078）
func _handle_vip_weekly_claimed(payload: Dictionary) -> void:
	emit_signal("vip_weekly_claimed", payload)

## 請求 VIP 狀態（DAY-078）
func request_vip_status() -> void:
	NetworkManager.send_message("get_vip_status", {})

## 領取 VIP 週獎勵（DAY-078）
func claim_vip_weekly() -> void:
	NetworkManager.send_message("claim_vip_weekly", {})

## 處理限時活動狀態更新（DAY-079）
func _handle_event_update(payload: Dictionary) -> void:
	emit_signal("event_updated", payload)

## 請求限時活動狀態（DAY-079）
func request_event_status() -> void:
	NetworkManager.send_message("get_event_status", {})

## 處理圖鑑狀態更新（DAY-081）
func _handle_codex_update(payload: Dictionary) -> void:
	emit_signal("codex_updated", payload)

## 處理圖鑑條目解鎖通知（DAY-081）
func _handle_codex_unlock(payload: Dictionary) -> void:
	emit_signal("codex_unlocked", payload)

## 處理全圖鑑完成通知（DAY-081）
func _handle_codex_complete(payload: Dictionary) -> void:
	emit_signal("codex_complete", payload)

## 請求圖鑑狀態（DAY-081）
func request_codex() -> void:
	NetworkManager.send_message("get_codex", {})

## 處理推薦碼資訊（DAY-082）
func _handle_referral_info(payload: Dictionary) -> void:
	emit_signal("referral_info_received", payload)

## 處理推薦碼使用成功（DAY-082）
func _handle_referral_success(payload: Dictionary) -> void:
	emit_signal("referral_success", payload)

## 處理推薦碼使用失敗（DAY-082）
func _handle_referral_error(payload: Dictionary) -> void:
	emit_signal("referral_error", payload)

## 請求推薦碼資訊（DAY-082）
func request_referral_info() -> void:
	NetworkManager.send_message("get_referral_info", {})

## 使用推薦碼（DAY-082）
func use_referral_code(code: String) -> void:
	NetworkManager.send_message("use_referral_code", {"code": code})

## 處理連擊狀態更新（DAY-083）
func _handle_streak_update(payload: Dictionary) -> void:
	emit_signal("streak_updated", payload)

## 處理連擊重置通知（DAY-083）
func _handle_streak_reset(payload: Dictionary) -> void:
	emit_signal("streak_reset", payload)

## 處理幸運轉盤觸發（DAY-084）
func _handle_wheel_trigger(payload: Dictionary) -> void:
	emit_signal("wheel_triggered", payload)

## 處理隱藏挑戰解鎖（DAY-085）
func _handle_challenge_unlocked(payload: Dictionary) -> void:
	emit_signal("challenge_unlocked", payload)

## 處理任務連續完成獎勵（DAY-086）
func _handle_mission_streak_bonus(payload: Dictionary) -> void:
	emit_signal("mission_streak_bonus", payload)

## 處理天氣狀態更新（DAY-087）
func _handle_weather_update(payload: Dictionary) -> void:
	emit_signal("weather_updated", payload)

## 處理連鎖爆炸（DAY-088）
func _handle_chain_explosion(payload: Dictionary) -> void:
	emit_signal("chain_explosion", payload)

## 處理特殊武器狀態更新（DAY-089）
func _handle_special_weapon_update(payload: Dictionary) -> void:
	emit_signal("special_weapon_updated", payload)

## 處理特殊武器發射廣播（DAY-089）
func _handle_special_weapon_fired(payload: Dictionary) -> void:
	emit_signal("special_weapon_fired", payload)

## 處理特殊武器自動充能完成（DAY-134）
func _handle_special_weapon_charged(payload: Dictionary) -> void:
	emit_signal("special_weapon_charged", payload)

## 處理追蹤飛彈命中結果（DAY-141）
func _handle_homing_missile_result(payload: Dictionary) -> void:
	emit_signal("homing_missile_result", payload)
	var killed: bool = payload.get("killed", false)
	var final_reward: int = payload.get("final_reward", 0)
	var multiplier: float = payload.get("multiplier", 0.0)
	if killed and final_reward > 0:
		print("[GameManager] Homing missile hit ×%.0f, reward=%d" % [multiplier, final_reward])

## 處理鑽頭龍蝦連帶效果（DAY-142）
func _handle_drill_lobster_chain(payload: Dictionary) -> void:
	emit_signal("drill_lobster_chain", payload)
	var phase: String = payload.get("phase", "")
	var total_reward: int = payload.get("total_reward", 0)
	if phase == "result" and total_reward > 0:
		print("[GameManager] Drill lobster chain result: reward=%d" % total_reward)

## 處理炸彈蟹連環爆炸（DAY-143）
func _handle_bomb_crab_chain(payload: Dictionary) -> void:
	emit_signal("bomb_crab_chain", payload)
	var phase: String = payload.get("phase", "")
	var total_reward: int = payload.get("total_reward", 0)
	if phase == "result" and total_reward > 0:
		print("[GameManager] Bomb crab chain result: reward=%d" % total_reward)

## 處理龍怒流星雨怒氣值更新（DAY-154）
func _handle_dragon_wrath_charge(payload: Dictionary) -> void:
	emit_signal("dragon_wrath_charge", payload)
	var progress: int = payload.get("progress", 0)
	var just_charged: bool = payload.get("just_charged", false)
	if just_charged:
		print("[GameManager] Dragon Wrath charged! progress=%d" % progress)

## 處理龍怒流星雨結果（DAY-154）
func _handle_dragon_wrath_result(payload: Dictionary) -> void:
	emit_signal("dragon_wrath_result", payload)
	var phase: String = payload.get("phase", "")
	var total_reward: int = payload.get("total_reward", 0)
	if phase == "result" and total_reward > 0:
		print("[GameManager] Dragon Wrath result: reward=%d" % total_reward)

## 處理魚雷爆炸結果（DAY-155）
func _handle_torpedo_result(payload: Dictionary) -> void:
	emit_signal("torpedo_result", payload)
	var phase: String = payload.get("phase", "")
	var total_reward: int = payload.get("total_reward", 0)
	if phase == "result" and total_reward > 0:
		print("[GameManager] Torpedo result: reward=%d" % total_reward)

## 處理軌道炮穿透結果（DAY-157）
func _handle_railgun_result(payload: Dictionary) -> void:
	emit_signal("railgun_result", payload)
	var phase: String = payload.get("phase", "")
	var total_reward: int = payload.get("total_reward", 0)
	if phase == "result" and total_reward > 0:
		print("[GameManager] Railgun result: reward=%d" % total_reward)

## 處理黑洞漩渦爆炸結果（DAY-166）
func _handle_black_hole_result(payload: Dictionary) -> void:
	emit_signal("black_hole_result", payload)
	var phase: String = payload.get("phase", "")
	var total_reward: int = payload.get("total_reward", 0)
	if phase == "result" and total_reward > 0:
		print("[GameManager] Black Hole result: reward=%d" % total_reward)

## 處理黃金輪盤螃蟹開始（DAY-167）
func _handle_roulette_crab_start(payload: Dictionary) -> void:
	emit_signal("roulette_crab_start", payload)
	var player_name: String = payload.get("player_name", "")
	print("[GameManager] Roulette Crab start: player=%s" % player_name)

## 處理黃金輪盤螃蟹結果（DAY-167）
func _handle_roulette_crab_result(payload: Dictionary) -> void:
	emit_signal("roulette_crab_result", payload)
	var wheel_result: float = payload.get("wheel_result", 0.0)
	var bonus_reward: int = payload.get("bonus_reward", 0)
	if bonus_reward > 0:
		print("[GameManager] Roulette Crab result: wheel=×%.0f reward=%d" % [wheel_result, bonus_reward])

## 處理獅子舞大獎爆發（DAY-168）
func _handle_lion_dance_burst(payload: Dictionary) -> void:
	emit_signal("lion_dance_burst", payload)
	var phase: String = payload.get("phase", "")
	var trigger_name: String = payload.get("trigger_name", "")
	var burst_mult: float = payload.get("burst_mult", 1.0)
	if phase == "burst_start":
		print("[GameManager] Lion Dance burst_start: player=%s mult=×%.0f" % [trigger_name, burst_mult])
	elif phase == "burst_end":
		var remaining: int = payload.get("remaining_targets", 0)
		print("[GameManager] Lion Dance burst_end: remaining=%d" % remaining)

## 處理漩渦魚群吸引（DAY-169）
func _handle_vortex_fish(payload: Dictionary) -> void:
	emit_signal("vortex_fish", payload)
	var phase: String = payload.get("phase", "")
	var trigger_name: String = payload.get("trigger_name", "")
	match phase:
		"vortex_start":
			var target_count: int = payload.get("target_count", 0)
			print("[GameManager] Vortex Fish vortex_start: player=%s targets=%d" % [trigger_name, target_count])
		"vortex_suck":
			var suck_index: int = payload.get("suck_index", 0)
			print("[GameManager] Vortex Fish vortex_suck: index=%d" % suck_index)
		"vortex_end":
			var killed_count: int = payload.get("killed_count", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Vortex Fish vortex_end: killed=%d reward=%d" % [killed_count, total_reward])

## 處理冰凍炸彈魚（DAY-170）
func _handle_freeze_bomb(payload: Dictionary) -> void:
	emit_signal("freeze_bomb", payload)
	var phase: String = payload.get("phase", "")
	var trigger_name: String = payload.get("trigger_name", "")
	var frozen_count: int = payload.get("frozen_count", 0)
	match phase:
		"freeze_start":
			print("[GameManager] Freeze Bomb freeze_start: player=%s frozen=%d" % [trigger_name, frozen_count])
		"freeze_end":
			print("[GameManager] Freeze Bomb freeze_end: frozen=%d" % frozen_count)

## 處理冰釣幸運輪盤（DAY-171）
func _handle_ice_fishing_wheel(payload: Dictionary) -> void:
	emit_signal("ice_fishing_wheel", payload)
	var phase: String = payload.get("phase", "")
	var player_name: String = payload.get("player_name", "")
	var multiplier: float = payload.get("multiplier", 1.0)
	match phase:
		"wheel_start":
			print("[GameManager] Ice Fishing wheel_start: player=%s" % player_name)
		"wheel_result":
			print("[GameManager] Ice Fishing wheel_result: mult=×%.0f" % multiplier)
		"mult_end":
			var kill_count: int = payload.get("kill_count", 0)
			var total_bonus: int = payload.get("total_bonus", 0)
			print("[GameManager] Ice Fishing mult_end: kills=%d bonus=%d" % [kill_count, total_bonus])

## 處理幸運彩蛋魚（DAY-172）
func _handle_lucky_egg_fish(payload: Dictionary) -> void:
	emit_signal("lucky_egg_fish", payload)
	var phase: String = payload.get("phase", "")
	var player_name: String = payload.get("player_name", "")
	var egg_count: int = payload.get("egg_count", 1)
	match phase:
		"egg_start":
			print("[GameManager] Lucky Egg Fish egg_start: player=%s eggs=%d" % [player_name, egg_count])
		"egg_open":
			var egg_index: int = payload.get("egg_index", 0)
			var reward_type: String = payload.get("egg_result", {}).get("reward_type", "")
			print("[GameManager] Lucky Egg Fish egg_open[%d]: type=%s" % [egg_index, reward_type])
		"egg_result":
			var total_coins: int = payload.get("total_coins", 0)
			var mult_count: int = payload.get("mult_count", 0)
			print("[GameManager] Lucky Egg Fish egg_result: eggs=%d coins=%d mult=%d" % [egg_count, total_coins, mult_count])
		"mult_end":
			print("[GameManager] Lucky Egg Fish mult_end")

## 處理彩虹幸運魚（DAY-173）
func _handle_rainbow_lucky_fish(payload: Dictionary) -> void:
	emit_signal("rainbow_lucky_fish", payload)
	var phase: String = payload.get("phase", "")
	var player_name: String = payload.get("player_name", "")
	match phase:
		"lucky_start":
			var duration_sec: int = payload.get("duration_sec", 10)
			var kill_boost: float = payload.get("kill_boost", 0.20)
			print("[GameManager] Rainbow Lucky Fish lucky_start: player=%s duration=%ds boost=+%.0f%%" % [player_name, duration_sec, kill_boost * 100])
		"lucky_end":
			print("[GameManager] Rainbow Lucky Fish lucky_end")

## 處理海葵觸手攻擊（DAY-174）
func _handle_sea_anemone(payload: Dictionary) -> void:
	emit_signal("sea_anemone", payload)
	var phase: String = payload.get("phase", "")
	var killer_name: String = payload.get("killer_name", "")
	match phase:
		"tentacle_start":
			print("[GameManager] Sea Anemone tentacle_start: player=%s" % killer_name)
		"tentacle_hit":
			var direction: int = payload.get("direction", 0)
			var is_kill: bool = payload.get("is_kill", false)
			var reward: int = payload.get("reward", 0)
			print("[GameManager] Sea Anemone tentacle_hit[%d]: kill=%s reward=%d" % [direction, is_kill, reward])
		"tentacle_result":
			var kill_count: int = payload.get("kill_count", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Sea Anemone tentacle_result: kills=%d reward=%d" % [kill_count, total_reward])

## 處理幸運骰子魚（DAY-175）
func _handle_lucky_dice_fish(payload: Dictionary) -> void:
	emit_signal("lucky_dice_fish", payload)
	var phase: String = payload.get("phase", "")
	var player_name: String = payload.get("player_name", "")
	match phase:
		"dice_start":
			print("[GameManager] Lucky Dice Fish dice_start: player=%s" % player_name)
		"dice_result":
			var die1: int = payload.get("die1", 1)
			var die2: int = payload.get("die2", 1)
			var sum: int = payload.get("sum", 2)
			var reward: int = payload.get("reward", 0)
			print("[GameManager] Lucky Dice Fish dice_result: %d+%d=%d reward=%d" % [die1, die2, sum, reward])
		"dice_jackpot":
			var reward: int = payload.get("reward", 0)
			print("[GameManager] Lucky Dice Fish dice_jackpot: reward=%d" % reward)

## 處理火焰風暴魚（DAY-176）
func _handle_fire_storm_fish(payload: Dictionary) -> void:
	emit_signal("fire_storm_fish", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"fire_start":
			var player_name: String = payload.get("player_name", "")
			var count: int = payload.get("target_count", 0)
			print("[GameManager] Fire Storm Fish fire_start: player=%s count=%d" % [player_name, count])
		"fire_end":
			var burned: int = payload.get("burned_count", 0)
			var reward: int = payload.get("total_reward", 0)
			print("[GameManager] Fire Storm Fish fire_end: burned=%d reward=%d" % [burned, reward])

## 處理黃金寶藏魚（DAY-177）
func _handle_golden_treasure_fish(payload: Dictionary) -> void:
	emit_signal("golden_treasure_fish", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"treasure_start":
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Golden Treasure Fish treasure_start: player=%s" % player_name)
		"treasure_open":
			var chest_id: int = payload.get("chest_id", 0)
			var reward_type: String = payload.get("reward_type", "coins")
			print("[GameManager] Golden Treasure Fish treasure_open: chest=%d type=%s" % [chest_id, reward_type])
		"treasure_end":
			print("[GameManager] Golden Treasure Fish treasure_end")

## send_golden_treasure_open — 發送開箱請求（DAY-177）
func send_golden_treasure_open(chest_id: int) -> void:
	NetworkManager.send("golden_treasure_open", {"chest_id": chest_id})

## 處理美人魚治癒（DAY-178）
func _handle_mermaid_healing(payload: Dictionary) -> void:
	emit_signal("mermaid_healing", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"heal_start":
			var heal_amount: int = payload.get("heal_amount", 0)
			print("[GameManager] Mermaid heal_start: amount=%d" % heal_amount)
		"luck_start":
			print("[GameManager] Mermaid luck_start: +20%% boost")
		"luck_end":
			print("[GameManager] Mermaid luck_end")

## 處理幸運草魚（DAY-179）
func _handle_lucky_clover_fish(payload: Dictionary) -> void:
	emit_signal("lucky_clover_fish", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"clover_start":
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Lucky Clover clover_start: player=%s" % player_name)
		"clover_gift":
			var gift_amount: int = payload.get("gift_amount", 0)
			print("[GameManager] Lucky Clover clover_gift: amount=%d" % gift_amount)
		"clover_end":
			print("[GameManager] Lucky Clover clover_end")

## 處理黃金海龜時間停止（DAY-159）
func _handle_golden_turtle_time_stop(payload: Dictionary) -> void:
	emit_signal("golden_turtle_time_stop", payload)
	var phase: String = payload.get("phase", "")
	print("[GameManager] Golden Turtle time stop: phase=%s" % phase)

## 處理幸運星魚全場倍率翻倍（DAY-160）
func _handle_lucky_star_fish(payload: Dictionary) -> void:
	emit_signal("lucky_star_fish", payload)
	var phase: String = payload.get("phase", "")
	print("[GameManager] Lucky Star Fish: phase=%s" % phase)

## 處理黃金鯊魚全服狂暴模式（DAY-161）
func _handle_golden_shark_berserk(payload: Dictionary) -> void:
	emit_signal("golden_shark_berserk", payload)
	var phase: String = payload.get("phase", "")
	print("[GameManager] Golden Shark Berserk: phase=%s" % phase)

## 處理金幣魚王即時獎勵（DAY-162）
func _handle_money_fish_reward(payload: Dictionary) -> void:
	emit_signal("money_fish_reward", payload)
	var reward: int = payload.get("instant_reward", 0)
	print("[GameManager] Money Fish King: instant_reward=%d" % reward)

## 處理船長魚全服競速模式（DAY-163）
func _handle_captain_fish_race(payload: Dictionary) -> void:
	emit_signal("captain_fish_race", payload)
	var phase: String = payload.get("phase", "")
	print("[GameManager] Captain Fish Race: phase=%s" % phase)

## 處理深淵巨鯨全服 Boss 挑戰（DAY-164）
func _handle_abyss_whale(payload: Dictionary) -> void:
	emit_signal("abyss_whale", payload)
	var phase: String = payload.get("phase", "")
	var current_hp: int = payload.get("current_hp", 0)
	var total_hp: int = payload.get("total_hp", 500)
	if phase == "whale_spawn":
		print("[GameManager] Abyss Whale spawned! HP=%d" % total_hp)
	elif phase == "whale_killed":
		var killer_name: String = payload.get("killer_name", "")
		print("[GameManager] Abyss Whale killed by %s!" % killer_name)
	elif phase == "whale_hp_update":
		pass # 頻繁更新，不 print

## 處理皇家閃電鰻持續連鎖電擊（DAY-156）
func _handle_royal_chain_lightning(payload: Dictionary) -> void:
	emit_signal("royal_chain_lightning", payload)
	var phase: String = payload.get("phase", "")
	var total_jumps: int = payload.get("total_jumps", 0)
	var total_reward: int = payload.get("total_reward", 0)
	if phase == "result" and total_reward > 0:
		print("[GameManager] Royal Chain Lightning result: jumps=%d reward=%d" % [total_jumps, total_reward])

## 處理失敗補償觸發（DAY-135）
func _handle_unlucky_bonus(payload: Dictionary) -> void:
	emit_signal("unlucky_bonus", payload)

## 處理神秘寶箱掉落通知（DAY-090）
func _handle_mystery_box_drop(payload: Dictionary) -> void:
	emit_signal("mystery_box_dropped", payload)

## 處理神秘寶箱狀態更新（DAY-090）
func _handle_mystery_box_update(payload: Dictionary) -> void:
	emit_signal("mystery_box_updated", payload)

## 處理神秘寶箱開箱結果（DAY-090）
func _handle_mystery_box_opened(payload: Dictionary) -> void:
	emit_signal("mystery_box_opened", payload)

## 處理房間列表（DAY-091）
func _handle_room_list(payload: Dictionary) -> void:
	emit_signal("room_list_received", payload)

## 處理房間切換成功（DAY-091）
func _handle_room_switched(payload: Dictionary) -> void:
	emit_signal("room_switched", payload)

## 處理房間操作失敗（DAY-091）
func _handle_room_error(payload: Dictionary) -> void:
	emit_signal("room_error", payload)

## 處理每日轉盤狀態（DAY-092）
func _handle_daily_spin_state(payload: Dictionary) -> void:
	emit_signal("daily_spin_state", payload)

## 處理每日轉盤結果（DAY-092）
func _handle_daily_spin_result(payload: Dictionary) -> void:
	emit_signal("daily_spin_result", payload)

## 處理登入里程碑達成通知（DAY-107）
func _handle_login_milestone(payload: Dictionary) -> void:	var days: int = payload.get("days", 0)
	var name: String = payload.get("name", "")
	print("[GameManager] Login milestone reached! Day=%d Name=%s" % [days, name])
	emit_signal("login_milestone_reached", payload)
	# 播放大獎音效
	if AudioManager != null:
		AudioManager.play_sfx("big_win")

## 處理登入進度回應（DAY-107）
func _handle_login_progress(payload: Dictionary) -> void:
	var streak: int = payload.get("current_streak", 0)
	print("[GameManager] Login progress: streak=%d" % streak)
	emit_signal("login_progress_received", payload)

## 請求登入進度（DAY-107）
func request_login_progress() -> void:
	NetworkManager.send_message("get_login_progress", {})

## 發送訊息（通用）
func send_message(type: String, payload: Dictionary) -> void:
	NetworkManager.send_message(type, payload)

## 處理超級 Bonus 通知（DAY-108）
func _handle_super_bonus_ready(payload: Dictionary) -> void:
	var label: String = payload.get("label", "SUPER BONUS!")
	var mult: float = payload.get("mult_bonus", 1.5)
	var combo: int = payload.get("combo_count", 3)
	print("[GameManager] %s combo=%d mult=%.1fx" % [label, combo, mult])
	emit_signal("super_bonus_triggered", payload)
	# 播放大獎音效
	if AudioManager != null:
		AudioManager.play_sfx("big_win")

## 發送訊息（通用）
func send_message(type: String, payload: Dictionary) -> void:
	NetworkManager.send_message(type, payload)

# ---- 賽季節日活動系統（DAY-109）----
signal festival_updated(festival_data: Dictionary)      # 節日狀態更新
signal festival_task_ready_signal(task_id: String)      # 節日任務可領取
signal festival_task_claimed_signal(task_data: Dictionary) # 節日任務獎勵領取成功
signal festival_title_earned_signal(title_data: Dictionary) # 節日稱號獲得
signal festival_error_signal(error_data: Dictionary)    # 節日操作失敗

## 處理節日狀態更新（DAY-109）
func _handle_festival_update(payload: Dictionary) -> void:
	var festival_type: String = payload.get("type", "none")
	var is_active: bool = payload.get("is_active", false)
	if is_active:
		print("[GameManager] Festival active: %s" % festival_type)
	emit_signal("festival_updated", payload)

## 處理節日任務可領取通知（DAY-109）
func _handle_festival_task_ready(payload: Dictionary) -> void:
	var task_id: String = payload.get("task_id", "")
	print("[GameManager] Festival task ready: %s" % task_id)
	emit_signal("festival_task_ready_signal", task_id)
	# 播放任務完成音效
	if AudioManager != null:
		AudioManager.play_sfx("bonus_ready")

## 處理節日任務獎勵領取成功（DAY-109）
func _handle_festival_task_claimed(payload: Dictionary) -> void:
	var task_id: String = payload.get("task_id", "")
	var coins: int = payload.get("reward_coins", 0)
	print("[GameManager] Festival task claimed: %s reward=%d" % [task_id, coins])
	emit_signal("festival_task_claimed_signal", payload)

## 處理節日稱號獲得通知（DAY-109）
func _handle_festival_title_earned(payload: Dictionary) -> void:
	var title_name: String = payload.get("title_name", "")
	print("[GameManager] Festival title earned: %s" % title_name)
	emit_signal("festival_title_earned_signal", payload)
	# 播放大獎音效
	if AudioManager != null:
		AudioManager.play_sfx("big_win")

## 處理節日操作失敗（DAY-109）
func _handle_festival_error(payload: Dictionary) -> void:
	emit_signal("festival_error_signal", payload)

## 請求節日狀態（DAY-109）
func request_festival() -> void:
	NetworkManager.send_message("get_festival", {})

## 發送領取節日任務獎勵請求（DAY-109）
func send_claim_festival_task(task_id: String) -> void:
	NetworkManager.send_message("claim_festival_task", {"task_id": task_id})

# ---- 名人堂系統（DAY-110）----

## 處理名人堂更新（DAY-110）
func _handle_hall_of_fame_update(payload: Dictionary) -> void:
	emit_signal("hall_of_fame_updated", payload)

## 處理新記錄誕生廣播（DAY-110）
func _handle_hall_of_fame_new_record(payload: Dictionary) -> void:
	var entry: Dictionary = payload.get("entry", {})
	var holder: String = entry.get("display_name", "")
	var label: String = entry.get("record_label", "")
	print("[GameManager] Hall of Fame NEW RECORD! %s by %s" % [label, holder])
	emit_signal("hall_of_fame_new_record", payload)
	# 播放大獎音效
	if AudioManager != null:
		AudioManager.play_sfx("big_win")

## 請求名人堂資料（DAY-110）
func request_hall_of_fame() -> void:
	NetworkManager.send_message("get_hall_of_fame", {})

# ---- 智慧推薦系統（DAY-110）----

## 處理推薦結果（DAY-110）
func _handle_recommendations(payload: Dictionary) -> void:
	var count: int = payload.get("recommendations", []).size()
	print("[GameManager] Recommendations received: %d" % count)
	emit_signal("recommendations_received", payload)

## 請求智慧推薦（DAY-110）
func request_recommendations() -> void:
	NetworkManager.send_message("get_recommendations", {})

## 發送投注等級切換（供推薦面板使用）
func send_bet_change(bet_level: int) -> void:
	NetworkManager.send_message("bet_change", {"bet_level": bet_level})

## 處理雙層倍率輪盤開始（DAY-113）
func _handle_roulette_start(payload: Dictionary) -> void:
	var player_id: String = payload.get("player_id", "")
	var player_name: String = payload.get("player_name", "")
	var target_name: String = payload.get("target_name", "")
	var is_self: bool = (player_id == get_player_id())
	payload["is_self"] = is_self
	print("[GameManager] Roulette started by %s on %s (self=%s)" % [player_name, target_name, is_self])
	emit_signal("roulette_started", payload)

## 處理雙層倍率輪盤結果（DAY-113）
func _handle_roulette_result(payload: Dictionary) -> void:
	var player_id: String = payload.get("player_id", "")
	var final_mult: float = payload.get("final_mult", 1.0)
	var final_reward: int = payload.get("final_reward", 0)
	var is_jackpot: bool = payload.get("is_jackpot", false)
	var is_self: bool = (player_id == get_player_id())
	payload["is_self"] = is_self
	if is_jackpot:
		print("[GameManager] Roulette JACKPOT! %.0fx reward=%d" % [final_mult, final_reward])
		if AudioManager != null:
			AudioManager.play_sfx("big_win")
	elif payload.get("is_mega_win", false):
		print("[GameManager] Roulette MEGA WIN! %.0fx reward=%d" % [final_mult, final_reward])
	emit_signal("roulette_result", payload)

# ---- 競速獵殺系統（DAY-136）----

## 處理競速獵殺開始廣播（DAY-136）
func _handle_speed_race_start(payload: Dictionary) -> void:
	var target_name: String = payload.get("target_name", "目標")
	var target_mult: float = payload.get("target_mult", 0.0)
	print("[GameManager] Speed Race started: %s (×%.0f)" % [target_name, target_mult])
	emit_signal("speed_race_started", payload)

## 處理競速獵殺結束廣播（DAY-136）
func _handle_speed_race_end(payload: Dictionary) -> void:
	var winner_name: String = payload.get("winner_name", "")
	var bonus_mult: float = payload.get("bonus_mult", 3.0)
	print("[GameManager] Speed Race ended! Winner: %s (×%.1f)" % [winner_name, bonus_mult])
	emit_signal("speed_race_ended", payload)
	if AudioManager != null:
		AudioManager.play_sfx("big_win")

## 處理競速獵殺取消廣播（DAY-136）
func _handle_speed_race_cancel(payload: Dictionary) -> void:
	var target_name: String = payload.get("target_name", "")
	print("[GameManager] Speed Race cancelled: %s" % target_name)
	emit_signal("speed_race_cancelled", payload)

## 處理競速個人結果（DAY-136）
func _handle_speed_race_result(payload: Dictionary) -> void:
	var rank: int = payload.get("rank", 0)
	var bonus_mult: float = payload.get("bonus_mult", 1.0)
	print("[GameManager] Speed Race result: rank=%d bonus=×%.1f" % [rank, bonus_mult])
	emit_signal("speed_race_result", payload)
	if rank == 1 and AudioManager != null:
		AudioManager.play_sfx("big_win")

# ---- 全服目標懸賞系統（DAY-137）----

## 處理懸賞發布廣播（DAY-137）
func _handle_bounty_posted(payload: Dictionary) -> void:
	var poster_name: String = payload.get("poster_name", "")
	var target_name: String = payload.get("target_name", "")
	var amount: int = payload.get("amount", 0)
	print("[GameManager] Bounty posted by %s on %s for %d coins" % [poster_name, target_name, amount])
	emit_signal("bounty_posted", payload)

## 處理懸賞個人領取通知（DAY-137）
func _handle_bounty_claimed(payload: Dictionary) -> void:
	var total_amount: int = payload.get("total_amount", 0)
	print("[GameManager] Bounty claimed! Total: %d coins" % total_amount)
	emit_signal("bounty_claimed", payload)
	if AudioManager != null:
		AudioManager.play_sfx("coin_drop")

## 處理懸賞目標擊破廣播（DAY-137）
func _handle_bounty_killed(payload: Dictionary) -> void:
	var killer_name: String = payload.get("killer_name", "")
	var total_amount: int = payload.get("total_amount", 0)
	print("[GameManager] Bounty killed by %s! Total: %d coins" % [killer_name, total_amount])
	emit_signal("bounty_killed", payload)

## 處理懸賞過期通知（DAY-137）
func _handle_bounty_expired(payload: Dictionary) -> void:
	var target_name: String = payload.get("target_name", "")
	print("[GameManager] Bounty expired: %s" % target_name)
	emit_signal("bounty_expired", payload)

## 對目標下懸賞（DAY-137）
func post_bounty(target_instance_id: String, amount: int) -> void:
	NetworkManager.send_message("post_bounty", {
		"target_instance_id": target_instance_id,
		"amount": amount
	})

## 查詢懸賞列表（DAY-137）
func request_bounties() -> void:
	NetworkManager.send_message("get_bounties", {})

# ---- 全服倍率風暴系統（DAY-138）----

## 處理倍率風暴開始廣播（DAY-138）
func _handle_mult_storm_start(payload: Dictionary) -> void:
	var tier_name: String = payload.get("tier_name", "⚡ 倍率風暴")
	var mult_boost: float = payload.get("mult_boost", 2.0)
	print("[GameManager] Mult Storm started: %s ×%.0f" % [tier_name, mult_boost])
	emit_signal("mult_storm_started", payload)
	if AudioManager != null:
		AudioManager.play_sfx("bonus_ready")

## 處理倍率風暴結束廣播（DAY-138）
func _handle_mult_storm_end(payload: Dictionary) -> void:
	print("[GameManager] Mult Storm ended")
	emit_signal("mult_storm_ended", payload)

## 處理雙環輪盤開始（DAY-139）
func _handle_dual_roulette_start(payload: Dictionary) -> void:
	var target_mult: float = payload.get("target_mult", 30.0)
	var base_reward: int = payload.get("base_reward", 0)
	print("[GameManager] Dual Roulette started: targetMult=%.0f, baseReward=%d" % [target_mult, base_reward])
	emit_signal("dual_roulette_started", payload)
	if AudioManager != null:
		AudioManager.play_sfx("bonus_ready")

## 處理雙環輪盤結果（DAY-139）
func _handle_dual_roulette_result(payload: Dictionary) -> void:
	var combined: float = payload.get("combined", 1.0)
	var bonus_reward: int = payload.get("bonus_reward", 0)
	print("[GameManager] Dual Roulette result: combined=%.0fx, bonus=%d" % [combined, bonus_reward])
	emit_signal("dual_roulette_result", payload)
	if combined >= 50.0 and AudioManager != null:
		AudioManager.play_sfx("big_win")

## 發送停止雙環輪盤請求（DAY-139）
func send_dual_roulette_stop() -> void:
	if NetworkManager != null:
		NetworkManager.send_message({"type": "dual_roulette_stop", "payload": {}})

## 處理 Mega Catch 事件開始（DAY-140）
func _handle_mega_catch_start(payload: Dictionary) -> void:
	var tier_name: String = payload.get("tier_name", "🎣 大豐收")
	var reward_boost: float = payload.get("reward_boost", 1.5)
	print("[GameManager] Mega Catch started: %s ×%.0f" % [tier_name, reward_boost])
	emit_signal("mega_catch_started", payload)
	if AudioManager != null:
		AudioManager.play_sfx("bonus_ready")

## 處理 Mega Catch 事件結束（DAY-140）
func _handle_mega_catch_end(payload: Dictionary) -> void:
	print("[GameManager] Mega Catch ended")
	emit_signal("mega_catch_ended", payload)

## 處理千龍王輪盤開始（DAY-148）
func _handle_chainlong_wheel_start(payload: Dictionary) -> void:
	var killer_name: String = payload.get("killer_name", "")
	var target_mult: float = payload.get("target_mult", 500.0)
	var base_reward: int = payload.get("base_reward", 0)
	print("[GameManager] ChainLong Wheel started: killer=%s, targetMult=%.0f, baseReward=%d" % [killer_name, target_mult, base_reward])
	emit_signal("chainlong_wheel_start", payload)
	if AudioManager != null:
		AudioManager.play_sfx("boss_enter")

## 處理千龍王輪盤結果（DAY-148）
func _handle_chainlong_wheel_result(payload: Dictionary) -> void:
	var combined: float = payload.get("combined", 1.0)
	var bonus_reward: int = payload.get("bonus_reward", 0)
	var is_mega_win: bool = payload.get("is_mega_win", false)
	print("[GameManager] ChainLong Wheel result: combined=%.0fx, bonus=%d, mega=%s" % [combined, bonus_reward, str(is_mega_win)])
	emit_signal("chainlong_wheel_result", payload)
	if is_mega_win and AudioManager != null:
		AudioManager.play_sfx("big_win")

## 發送停止千龍王輪盤請求（DAY-148）
func send_chainlong_wheel_stop() -> void:
	if NetworkManager != null:
		NetworkManager.send_message({"type": "chainlong_wheel_stop", "payload": {}})

## 處理彩虹鯊魚爆發（DAY-180）
func _handle_rainbow_shark_burst(payload: Dictionary) -> void:
	emit_signal("rainbow_shark_burst", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"burst_start":
			var trigger_name: String = payload.get("trigger_name", "")
			var marked_count: int = payload.get("marked_targets", []).size()
			print("[GameManager] Rainbow Shark burst_start: player=%s, targets=%d" % [trigger_name, marked_count])
		"burst_end":
			print("[GameManager] Rainbow Shark burst_end")

## 處理雷霆鯊魚連鎖閃電（DAY-181）
func _handle_thunder_shark_chain(payload: Dictionary) -> void:
	emit_signal("thunder_shark_chain", payload)
	var phase: String = payload.get("phase", "")
	if phase == "chain_start":
		var killer_name: String = payload.get("killer_name", "")
		print("[GameManager] Thunder Shark chain_start: killer=%s" % killer_name)
	elif phase == "result":
		var total_jumps: int = payload.get("total_jumps", 0)
		var total_kills: int = payload.get("total_kills", 0)
		print("[GameManager] Thunder Shark result: jumps=%d kills=%d" % [total_jumps, total_kills])

## 處理吸血鬼魚累積倍率（DAY-182）
func _handle_vampire_fish(payload: Dictionary) -> void:
	emit_signal("vampire_fish", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"vampire_start":
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Vampire Fish vampire_start: player=%s" % player_name)
		"mult_update":
			var current_mult: float = payload.get("current_mult", 1.0)
			print("[GameManager] Vampire Fish mult_update: mult=%.1f" % current_mult)
		"vampire_end":
			var final_mult: float = payload.get("current_mult", 1.0)
			var kill_count: int = payload.get("kill_count", 0)
			print("[GameManager] Vampire Fish vampire_end: finalMult=%.1f kills=%d" % [final_mult, kill_count])

## 處理閃電魚自動連鎖（DAY-183）
func _handle_lightning_auto_chain(payload: Dictionary) -> void:
	emit_signal("lightning_auto_chain", payload)
	var phase: String = payload.get("phase", "")
	if phase == "chain_start":
		var killer_name: String = payload.get("killer_name", "")
		print("[GameManager] Lightning Auto Chain start: killer=%s" % killer_name)
	elif phase == "result":
		var total_attacks: int = payload.get("total_attacks", 0)
		var total_kills: int = payload.get("total_kills", 0)
		print("[GameManager] Lightning Auto Chain result: attacks=%d kills=%d" % [total_attacks, total_kills])

## 處理隕石魚隕石雨（DAY-184）
func _handle_meteor_fish(payload: Dictionary) -> void:
	emit_signal("meteor_fish", payload)
	var phase: String = payload.get("phase", "")
	if phase == "meteor_start":
		var killer_name: String = payload.get("killer_name", "")
		var meteor_count: int = payload.get("meteor_count", 0)
		print("[GameManager] Meteor Fish shower start: killer=%s meteors=%d" % [killer_name, meteor_count])
	elif phase == "meteor_result":
		var total_kills: int = payload.get("total_kills", 0)
		var total_reward: int = payload.get("total_reward", 0)
		print("[GameManager] Meteor Fish result: kills=%d reward=%d" % [total_kills, total_reward])

## 處理鳳凰魚涅槃重生（DAY-185）
func _handle_phoenix_fish(payload: Dictionary) -> void:
	emit_signal("phoenix_fish", payload)
	var phase: String = payload.get("phase", "")
	if phase == "phoenix_explode":
		var killer_name: String = payload.get("killer_name", "")
		print("[GameManager] Phoenix Fish explode: killer=%s" % killer_name)
	elif phase == "phoenix_rebirth":
		var total_kills: int = payload.get("total_kills", 0)
		var boost_sec: int = payload.get("boost_sec", 30)
		print("[GameManager] Phoenix Fish rebirth: kills=%d boost=%ds" % [total_kills, boost_sec])
	elif phase == "rebirth_end":
		print("[GameManager] Phoenix Fish rebirth ended")

## 處理龍龜不死 Boss（DAY-186）
func _handle_dragon_turtle(payload: Dictionary) -> void:
	emit_signal("dragon_turtle", payload)
	var phase: String = payload.get("phase", "")
	if phase == "turtle_appear":
		var instance_id: String = payload.get("instance_id", "")
		print("[GameManager] Dragon Turtle appeared: instance=%s" % instance_id)
	elif phase == "turtle_hit":
		var hitter_name: String = payload.get("hitter_name", "")
		var hit_mult: int = payload.get("hit_mult", 0)
		print("[GameManager] Dragon Turtle hit: hitter=%s mult=%d" % [hitter_name, hit_mult])
	elif phase == "my_hit":
		var hit_reward: int = payload.get("hit_reward", 0)
		var hit_mult: int = payload.get("hit_mult", 0)
		print("[GameManager] Dragon Turtle my_hit: reward=%d mult=%d" % [hit_reward, hit_mult])
	elif phase == "turtle_leave":
		var total_hits: int = payload.get("total_hits", 0)
		var total_reward: int = payload.get("total_reward", 0)
		print("[GameManager] Dragon Turtle left: hits=%d reward=%d" % [total_hits, total_reward])

## 處理連鎖爆炸魚（DAY-187）
func _handle_chain_bomb(payload: Dictionary) -> void:
	emit_signal("chain_bomb", payload)
	var phase: String = payload.get("phase", "")
	if phase == "chain_start":
		var killer_name: String = payload.get("killer_name", "")
		print("[GameManager] Chain Bomb started: killer=%s" % killer_name)
	elif phase == "chain_explode":
		var chain_depth: int = payload.get("chain_depth", 0)
		var kill_count: int = payload.get("kill_count", 0)
		print("[GameManager] Chain Bomb explode: depth=%d kills=%d" % [chain_depth, kill_count])
	elif phase == "chain_result":
		var total_kills: int = payload.get("total_kills", 0)
		var total_reward: int = payload.get("total_reward", 0)
		var chain_depth: int = payload.get("chain_depth", 0)
		print("[GameManager] Chain Bomb result: depth=%d kills=%d reward=%d" % [chain_depth, total_kills, total_reward])

## 處理巨型鱷魚獵食（DAY-188）
func _handle_crocodile_hunter(payload: Dictionary) -> void:
	emit_signal("crocodile_hunter", payload)
	var phase: String = payload.get("phase", "")
	if phase == "croc_appear":
		var max_hunts: int = payload.get("max_hunts", 8)
		print("[GameManager] Crocodile Hunter appeared: max_hunts=%d" % max_hunts)
	elif phase == "croc_hunt":
		var hunt_index: int = payload.get("hunt_index", 0)
		var target_name: String = payload.get("target_name", "")
		var hunt_reward: int = payload.get("hunt_reward", 0)
		var total_pool: int = payload.get("total_pool", 0)
		print("[GameManager] Croc hunt #%d: target=%s reward=%d pool=%d" % [hunt_index, target_name, hunt_reward, total_pool])
	elif phase == "croc_killed":
		var killer_name: String = payload.get("killer_name", "")
		var total_reward: int = payload.get("total_reward", 0)
		print("[GameManager] Croc killed by %s: total_reward=%d" % [killer_name, total_reward])
	elif phase == "croc_leave":
		var hunt_count: int = payload.get("hunt_count", 0)
		var total_pool: int = payload.get("total_pool", 0)
		print("[GameManager] Croc left: hunts=%d pool=%d" % [hunt_count, total_pool])

## 處理時間炸彈魚（DAY-189）
func _handle_time_bomb_fish(payload: Dictionary) -> void:
	emit_signal("time_bomb_fish", payload)
	var phase: String = payload.get("phase", "")
	if phase == "bomb_appear":
		var countdown: int = payload.get("countdown", 10)
		print("[GameManager] Time Bomb appeared: countdown=%d" % countdown)
	elif phase == "bomb_tick":
		var remaining: int = payload.get("countdown", 0)
		print("[GameManager] Time Bomb tick: remaining=%d" % remaining)
	elif phase == "bomb_defused":
		var killer_name: String = payload.get("killer_name", "")
		var bonus_pct: int = payload.get("bonus_pct", 25)
		print("[GameManager] Time Bomb defused by %s: bonus=+%d%%" % [killer_name, bonus_pct])
	elif phase == "bomb_explode":
		print("[GameManager] Time Bomb EXPLODED!")
	elif phase == "bomb_result":
		var kill_count: int = payload.get("kill_count", 0)
		var total_reward: int = payload.get("total_reward", 0)
		print("[GameManager] Time Bomb result: kills=%d reward=%d" % [kill_count, total_reward])

## 處理三重幸運魚（DAY-190）
func _handle_triple_lucky_fish(payload: Dictionary) -> void:
	emit_signal("triple_lucky_fish", payload)
	var phase: String = payload.get("phase", "")
	if phase == "triple_start":
		var coin_reward: int = payload.get("coin_reward", 0)
		var weapon_charged: String = payload.get("weapon_charged", "")
		print("[GameManager] Triple Lucky Fish triggered: coin=%d weapon=%s" % [coin_reward, weapon_charged])
	elif phase == "triple_broadcast":
		var player_name: String = payload.get("player_name", "")
		print("[GameManager] Triple Lucky Fish broadcast: player=%s" % player_name)
	elif phase == "mult_end":
		print("[GameManager] Triple Lucky Fish mult bonus ended")

## 處理魚群驚嚇連帶（DAY-191）
func _handle_school_panic(payload: Dictionary) -> void:
	emit_signal("school_panic", payload)
	var phase: String = payload.get("phase", "")
	if phase == "panic_start":
		var target_count: int = payload.get("target_count", 0)
		var killer_name: String = payload.get("killer_name", "")
		print("[GameManager] School Panic triggered by %s: %d targets HP halved" % [killer_name, target_count])
	elif phase == "panic_end":
		print("[GameManager] School Panic ended")

## 處理搖滾骷髏演唱會（DAY-192）
func _handle_rock_skeleton_concert(payload: Dictionary) -> void:
	emit_signal("rock_skeleton_concert", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"concert_start":
			var killer_name: String = payload.get("killer_name", "")
			print("[GameManager] Rock Skeleton Concert started by %s" % killer_name)
		"awakening":
			var awakened_count: int = payload.get("awakened_count", 0)
			print("[GameManager] Rock Skeleton Awakening! %d targets HP reduced 70%%" % awakened_count)
		"encore_start":
			var total_kills: int = payload.get("total_kills", 0)
			print("[GameManager] Rock Skeleton Encore! kills=%d +30%% bonus activated" % total_kills)
		"concert_end":
			var total_kills: int = payload.get("total_kills", 0)
			print("[GameManager] Rock Skeleton Concert ended. kills=%d" % total_kills)

## 處理電流水母電流網路（DAY-193）
func _handle_electric_jellyfish(payload: Dictionary) -> void:
	emit_signal("electric_jellyfish", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"network_start":
			var killer_name: String = payload.get("killer_name", "")
			var link_count: int = payload.get("link_count", 0)
			print("[GameManager] Electric Jellyfish network started by %s: %d links" % [killer_name, link_count])
		"network_result":
			var total_kills: int = payload.get("total_kills", 0)
			var link_count: int = payload.get("link_count", 0)
			print("[GameManager] Electric Jellyfish network ended. links=%d kills=%d" % [link_count, total_kills])

## 處理長龍王雙環輪盤（DAY-194）
func _handle_chainlong_king(payload: Dictionary) -> void:
	emit_signal("chainlong_king", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"roulette_start":
			var instance_id: String = payload.get("instance_id", "")
			print("[GameManager] ChainLong King roulette started: %s" % instance_id)
		"mega_win":
			var reward: int = payload.get("reward", 0)
			print("[GameManager] ChainLong King MEGA WIN! reward=%d" % reward)
		"result":
			var total_mult: int = payload.get("total_mult", 0)
			var reward: int = payload.get("reward", 0)
			print("[GameManager] ChainLong King result: %dx reward=%d" % [total_mult, reward])

## 發送長龍王停止輪盤訊息（DAY-194）
func send_chainlong_king_stop(instance_id: String) -> void:
	_send_message("chainlong_king_stop", {"instance_id": instance_id})

## 處理巨型鮟鱇魚電擊寶箱（DAY-196）
func _handle_anglerfish_electric(payload: Dictionary) -> void:
	emit_signal("anglerfish_electric", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"anglerfish_appear":
			var instance_id: String = payload.get("instance_id", "")
			print("[GameManager] Anglerfish Electric appeared: %s" % instance_id)
		"anglerfish_killed":
			var killer_name: String = payload.get("killer_name", "")
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Anglerfish killed by %s, reward=%d" % [killer_name, total_reward])
		"super_zap_result":
			var super_kills: int = payload.get("super_kills", 0)
			var super_reward: int = payload.get("super_reward", 0)
			print("[GameManager] Super Zap result: kills=%d reward=%d" % [super_kills, super_reward])

## 處理神秘龍魚八波攻擊（DAY-197）
func _handle_mystic_dragon(payload: Dictionary) -> void:
	emit_signal("mystic_dragon", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"dragon_start":
			var killer_name: String = payload.get("killer_name", "")
			print("[GameManager] Mystic Dragon 8-wave attack started by %s" % killer_name)
		"dragon_result":
			var total_kills: int = payload.get("total_kills", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Mystic Dragon result: kills=%d reward=%d" % [total_kills, total_reward])

## 處理幽靈魚分身（DAY-198）
func _handle_ghost_fish(payload: Dictionary) -> void:
	emit_signal("ghost_fish", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"ghost_appear":
			var clone_count: int = payload.get("clone_count", 2)
			print("[GameManager] Ghost Fish appeared with %d clones" % clone_count)
		"real_found":
			var killer_name: String = payload.get("killer_name", "")
			print("[GameManager] Ghost Fish real body found by %s" % killer_name)
		"ghost_explode":
			var explode_kills: int = payload.get("explode_kills", 0)
			var explode_reward: int = payload.get("explode_reward", 0)
			print("[GameManager] Ghost Fish explode: kills=%d reward=%d" % [explode_kills, explode_reward])

## 處理雷霆龍蝦免費射擊（DAY-199）
func _handle_thunderbolt_lobster(payload: Dictionary) -> void:
	emit_signal("thunderbolt_lobster", payload)
	var event: String = payload.get("event", "")
	match event:
		"turret_start":
			var killer_name: String = payload.get("killer_name", "")
			print("[GameManager] Thunderbolt Lobster turret started by %s" % killer_name)
		"turret_end":
			var kill_count: int = payload.get("kill_count", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Thunderbolt Lobster turret ended: kills=%d reward=%d" % [kill_count, total_reward])

## 處理冰鳳凰覺醒 BOSS（DAY-200）
func _handle_ice_phoenix(payload: Dictionary) -> void:
	emit_signal("ice_phoenix", payload)
	var event: String = payload.get("event", "")
	match event:
		"awaken_start":
			var killer_name: String = payload.get("killer_name", "")
			print("[GameManager] Ice Phoenix awakened by %s" % killer_name)
		"awaken_result":
			var total_reward: int = payload.get("total_reward", 0)
			var has_frost: bool = payload.get("has_frost", false)
			print("[GameManager] Ice Phoenix result: total=%d frost=%s" % [total_reward, str(has_frost)])

## 處理連環炸彈蟹（DAY-201）
func _handle_serial_bomb_crab(payload: Dictionary) -> void:
	emit_signal("serial_bomb_crab", payload)
	var event: String = payload.get("event", "")
	match event:
		"bomb_start":
			var bomb_count: int = payload.get("bomb_count", 3)
			print("[GameManager] Serial Bomb Crab: %d bombs" % bomb_count)
		"bomb_result":
			var total_kills: int = payload.get("total_kills", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Serial Bomb Crab result: kills=%d reward=%d" % [total_kills, total_reward])

## 處理深淵漩渦魚（DAY-202）
func _handle_abyss_vortex(payload: Dictionary) -> void:
	emit_signal("abyss_vortex", payload)
	var event: String = payload.get("event", "")
	match event:
		"vortex_start":
			var killer_name: String = payload.get("killer_name", "")
			print("[GameManager] Abyss Vortex started by %s" % killer_name)
		"vortex_result":
			var total_kills: int = payload.get("total_kills", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Abyss Vortex result: kills=%d reward=%d" % [total_kills, total_reward])

## 處理座頭鯨覺醒（DAY-203）
func _handle_humpback_whale(payload: Dictionary) -> void:
	emit_signal("humpback_whale", payload)
	var event: String = payload.get("event", "")
	match event:
		"awaken_start":
			var killer_name: String = payload.get("killer_name", "")
			print("[GameManager] Humpback Whale awaken by %s" % killer_name)
		"awaken_result":
			var total_kills: int = payload.get("total_kills", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Humpback Whale result: kills=%d reward=%d" % [total_kills, total_reward])

## 處理自由旋轉魚免費射擊（DAY-204）
func _handle_free_spin_fish(payload: Dictionary) -> void:
	emit_signal("free_spin_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"free_spin_start":
			var player_id: String = payload.get("player_id", "")
			print("[GameManager] Free Spin Fish started for player=%s" % player_id)
		"free_spin_end":
			var kill_count: int = payload.get("kill_count", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Free Spin Fish ended: kills=%d reward=%d" % [kill_count, total_reward])

## 處理獎池龍 Jackpot 抽獎（DAY-205）
func _handle_jackpot_dragon(payload: Dictionary) -> void:
	emit_signal("jackpot_dragon", payload)
	var level: String = payload.get("level", "mini")
	var amount: int = payload.get("amount", 0)
	var player_name: String = payload.get("player_name", "")
	print("[GameManager] Jackpot Dragon: player=%s level=%s amount=%d" % [player_name, level, amount])

## 處理彗星魚連鎖爆炸（DAY-206）
func _handle_comet_fish(payload: Dictionary) -> void:
	emit_signal("comet_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"comet_appear":
			var trail_count: int = payload.get("trail_count", 7)
			print("[GameManager] Comet Fish appeared, trail=%d" % trail_count)
		"trail_blast":
			var blast_idx: int = payload.get("blast_index", 1)
			var kills: int = payload.get("kill_count", 0)
			print("[GameManager] Comet trail blast #%d: kills=%d" % [blast_idx, kills])
		"early_supernova":
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Comet early supernova by %s" % player_name)
		"supernova":
			var kills: int = payload.get("kill_count", 0)
			var reward: int = payload.get("reward", 0)
			print("[GameManager] Comet supernova: kills=%d reward=%d" % [kills, reward])

## 處理黃金波浪魚全場倍率衝擊（DAY-207）
func _handle_golden_wave_fish(payload: Dictionary) -> void:
	emit_signal("golden_wave_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"wave_start":
			var boost_mult: float = payload.get("boost_mult", 2.0)
			print("[GameManager] Golden Wave started, boost=×%.0f" % boost_mult)
		"wave_column":
			var col_index: int = payload.get("col_index", 0)
			var kills: int = payload.get("kill_count", 0)
			print("[GameManager] Golden Wave col #%d: kills=%d" % [col_index, kills])
		"boost_start":
			var total_kills: int = payload.get("total_kills", 0)
			var boost_sec: int = payload.get("boost_sec", 8)
			print("[GameManager] Golden Wave boost started: kills=%d, %ds" % [total_kills, boost_sec])
		"boost_end":
			print("[GameManager] Golden Wave boost ended")

## 處理深海龍王全服合力蓄力（DAY-208）
func _handle_dragon_king(payload: Dictionary) -> void:
	emit_signal("dragon_king", payload)
	var event: String = payload.get("event", "")
	match event:
		"charge_start":
			var charge_target: int = payload.get("charge_target", 20)
			var charge_sec: int = payload.get("charge_sec", 12)
			print("[GameManager] Dragon King charge started: target=%d, %ds" % [charge_target, charge_sec])
		"charge_progress":
			var current: int = payload.get("current", 0)
			var target: int = payload.get("charge_target", 20)
			print("[GameManager] Dragon King charge: %d/%d" % [current, target])
		"meteor_rain_start":
			var meteor_count: int = payload.get("meteor_count", 5)
			print("[GameManager] Dragon King FULL meteor rain: %d meteors" % meteor_count)
		"small_meteor_start":
			var meteor_count: int = payload.get("meteor_count", 3)
			print("[GameManager] Dragon King small meteor: %d meteors" % meteor_count)
		"meteor_hit":
			var kills: int = payload.get("kill_count", 0)
			var reward: int = payload.get("reward", 0)
			print("[GameManager] Dragon King meteor hit: kills=%d reward=%d" % [kills, reward])
		"meteor_rain_result":
			var total_kills: int = payload.get("total_kills", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Dragon King meteor rain result: kills=%d reward=%d" % [total_kills, total_reward])
		"small_meteor_result":
			var total_kills: int = payload.get("total_kills", 0)
			print("[GameManager] Dragon King small meteor result: kills=%d" % total_kills)

## 處理幸運金幣魚即時獎勵（DAY-209）
func _handle_fortune_coin_fish(payload: Dictionary) -> void:
	emit_signal("fortune_coin_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"coin_burst":
			var multiplier: int = payload.get("multiplier", 5)
			var reward: int = payload.get("reward", 0)
			print("[GameManager] Fortune Coin burst: ×%d reward=%d" % [multiplier, reward])
		"coin_broadcast":
			var player_name: String = payload.get("player_name", "")
			var label: String = payload.get("label", "")
			print("[GameManager] Fortune Coin broadcast: %s %s" % [player_name, label])
		"golden_burst_start":
			var affected: int = payload.get("affected_count", 0)
			var burst_sec: int = payload.get("burst_sec", 5)
			print("[GameManager] Fortune Coin Golden Burst: affected=%d, %ds" % [affected, burst_sec])
		"golden_burst_end":
			print("[GameManager] Fortune Coin Golden Burst ended")

## 處理幸運熱區魚空間策略（DAY-210）
func _handle_lucky_hot_zone(payload: Dictionary) -> void:
	emit_signal("lucky_hot_zone", payload)
	var event: String = payload.get("event", "")
	match event:
		"zone_start":
			var player_name: String = payload.get("player_name", "")
			var zone_x: float = payload.get("zone_x", 640.0)
			var zone_y: float = payload.get("zone_y", 360.0)
			var duration: int = payload.get("duration_sec", 8)
			print("[GameManager] Lucky Hot Zone start: %s at (%.0f,%.0f) %ds" % [player_name, zone_x, zone_y, duration])
		"zone_pulse":
			var affected: int = payload.get("affected_count", 0)
			var remaining: int = payload.get("remaining_sec", 0)
			print("[GameManager] Lucky Hot Zone pulse: affected=%d remaining=%ds" % [affected, remaining])
		"zone_blast":
			var killed: int = payload.get("killed_count", 0)
			var reward_per: int = payload.get("reward_per_player", 0)
			print("[GameManager] Lucky Hot Zone blast: killed=%d reward_per=%d" % [killed, reward_per])

## 處理幸運三叉魚互動三轉盤（DAY-211）
func _handle_lucky_trident(payload: Dictionary) -> void:
	emit_signal("lucky_trident", payload)
	var event: String = payload.get("event", "")
	match event:
		"trident_start":
			var timeout: int = payload.get("timeout_sec", 12)
			print("[GameManager] Lucky Trident start: timeout=%ds" % timeout)
		"trident_result":
			var a: String = payload.get("wheel_a_label", "")
			var b: String = payload.get("wheel_b_label", "")
			var c: String = payload.get("wheel_c_label", "")
			var coin: int = payload.get("coin_reward", 0)
			print("[GameManager] Lucky Trident result: %s + %s + %s coin=%d" % [a, b, c, coin])
		"trident_effect":
			var effect: String = payload.get("effect", "")
			print("[GameManager] Lucky Trident effect: %s" % effect)
		"trident_broadcast":
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Lucky Trident broadcast: %s" % player_name)

## 處理時間凍結魚系統（DAY-212）
func _handle_time_freeze_fish(payload: Dictionary) -> void:
	emit_signal("time_freeze_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"freeze_start":
			var duration: int = payload.get("duration_sec", 5)
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Time Freeze start: player=%s duration=%ds" % [player_name, duration])
		"freeze_end":
			var hit_count: int = payload.get("hit_count", 0)
			print("[GameManager] Time Freeze end: hit_count=%d" % hit_count)
		"thaw_blast":
			var killed: int = payload.get("killed_count", 0)
			var reward_per: int = payload.get("reward_per_player", 0)
			print("[GameManager] Time Freeze thaw blast: killed=%d reward_per=%d" % [killed, reward_per])

## 處理彩虹稜鏡魚系統（DAY-213）
func _handle_rainbow_prism(payload: Dictionary) -> void:
	emit_signal("rainbow_prism", payload)
	var event: String = payload.get("event", "")
	match event:
		"prism_start":
			var trigger_player: String = payload.get("trigger_player", "")
			var count: int = payload.get("colored_targets", []).size()
			print("[GameManager] Rainbow Prism start: player=%s colored=%d" % [trigger_player, count])
		"prism_blast":
			var kills: int = payload.get("blast_kills", 0)
			var reward: int = payload.get("blast_reward", 0)
			print("[GameManager] Rainbow Prism blast: kills=%d reward=%d" % [kills, reward])

## 處理黃金累積魚系統（DAY-214）
func _handle_golden_accumulator(payload: Dictionary) -> void:
	emit_signal("golden_accumulator", payload)
	var event: String = payload.get("event", "")
	match event:
		"accum_appear":
			var target: int = payload.get("accum_target", 20)
			print("[GameManager] Golden Accumulator appeared: target=%d" % target)
		"accum_progress":
			var count: int = payload.get("accum_count", 0)
			var tgt: int = payload.get("accum_target", 20)
			print("[GameManager] Golden Accumulator progress: %d/%d" % [count, tgt])
		"burst_start", "early_burst_start":
			var boost: float = payload.get("boost_mult", 2.0)
			var sec: int = payload.get("boost_sec", 8)
			print("[GameManager] Golden Accumulator burst: boost=%.1fx %ds" % [boost, sec])
		"burst_end":
			print("[GameManager] Golden Accumulator burst ended")

## 處理幸運鏡像魚系統（DAY-215）
func _handle_lucky_mirror_fish(payload: Dictionary) -> void:
	emit_signal("lucky_mirror_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"mirror_start":
			var mirrors: Array = payload.get("mirrors", [])
			var mult: float = payload.get("mult_boost", 1.5)
			print("[GameManager] Lucky Mirror Fish: %d mirrors created (x%.1f)" % [mirrors.size(), mult])
		"mirror_kill":
			var mirror_id: String = payload.get("mirror_id", "")
			print("[GameManager] Mirror killed: %s" % mirror_id)
		"mirror_blast":
			var blast_count: int = payload.get("blast_count", 0)
			print("[GameManager] Mirror blast: %d" % blast_count)
		"mirror_result":
			var killed: int = payload.get("killed_count", 0)
			var blasted: int = payload.get("blast_count", 0)
			var reward: int = payload.get("total_reward", 0)
			print("[GameManager] Mirror result: killed=%d blast=%d reward=%d" % [killed, blasted, reward])

## 處理詛咒毒魚系統（DAY-216）
func _handle_cursed_poison_fish(payload: Dictionary) -> void:
	emit_signal("cursed_poison_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"curse_start":
			var targets: Array = payload.get("cursed_targets", [])
			var mult: float = payload.get("curse_mult", 2.5)
			print("[GameManager] Cursed Poison Fish: %d targets cursed (x%.1f)" % [targets.size(), mult])
		"curse_kill":
			var instance_id: String = payload.get("instance_id", "")
			print("[GameManager] Cursed target killed: %s" % instance_id)
		"curse_escape":
			var penalty_sec: int = payload.get("penalty_sec", 5)
			print("[GameManager] Curse escape! Penalty x0.5 for %ds" % penalty_sec)
		"curse_cleanse":
			var player_name: String = payload.get("player_name", "")
			var reward: int = payload.get("cleanse_reward", 0)
			print("[GameManager] Curse cleansed by %s, reward=%d" % [player_name, reward])

## 處理鑽頭龍蝦穿透爆炸（DAY-195）
func _handle_drill_lobster(payload: Dictionary) -> void:
	emit_signal("drill_lobster", payload)
	var phase: String = payload.get("phase", "")
	match phase:
		"drill_start":
			var killer_name: String = payload.get("killer_name", "")
			print("[GameManager] Drill Lobster started by %s" % killer_name)
		"drill_result":
			var total_kills: int = payload.get("total_kills", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Drill Lobster result: kills=%d reward=%d" % [total_kills, total_reward])

## 處理幸運拍賣魚系統（DAY-217）
func _handle_lucky_auction_fish(payload: Dictionary) -> void:
	emit_signal("lucky_auction_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"auction_start":
			var duration: int = payload.get("duration_sec", 8)
			var control_mult: float = payload.get("control_mult", 0.85)
			print("[GameManager] Lucky Auction Fish: auction started (%ds, x%.2f)" % [duration, control_mult])
		"auction_bid":
			var player_name: String = payload.get("player_name", "")
			var bid_amount: int = payload.get("bid_amount", 0)
			var top_bidder: String = payload.get("top_bidder", "")
			print("[GameManager] Auction bid: %s bid=%d, top=%s" % [player_name, bid_amount, top_bidder])
		"auction_result":
			var winner: String = payload.get("winner_name", "")
			var winner_bid: int = payload.get("winner_bid", 0)
			print("[GameManager] Auction result: winner=%s bid=%d" % [winner, winner_bid])
		"auction_no_bid":
			print("[GameManager] Auction ended with no bids")
		"auction_fish_killed":
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Auction fish killed by %s" % player_name)
		"control_shot":
			var shot_reward: int = payload.get("shot_reward", 0)
			var shot_count: int = payload.get("shot_count", 0)
			print("[GameManager] Control shot: count=%d reward=%d" % [shot_count, shot_reward])
		"control_end":
			var player_name: String = payload.get("player_name", "")
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Control ended: player=%s total=%d" % [player_name, total_reward])

## 處理幸運進化魚系統（DAY-218）
func _handle_lucky_evolution_fish(payload: Dictionary) -> void:
	emit_signal("lucky_evolution_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"evolution_appear":
			print("[GameManager] Lucky Evolution Fish appeared")
		"evolution_hit":
			var hit_count: int = payload.get("hit_count", 0)
			var next_hit: int = payload.get("next_hit", 3)
			print("[GameManager] Evolution hit: %d/%d" % [hit_count, next_hit])
		"evolution_stage":
			var stage: int = payload.get("stage", 1)
			var mult: float = payload.get("mult_boost", 1.5)
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Evolution stage %d: player=%s mult=x%.1f" % [stage, player_name, mult])
		"evolution_burst", "evolution_kill_burst":
			var mult: float = payload.get("mult_boost", 4.0)
			var affected: int = payload.get("affected_count", 0)
			print("[GameManager] Evolution burst: mult=x%.1f affected=%d" % [mult, affected])
		"evolution_burst_end":
			print("[GameManager] Evolution burst ended")
		"evolution_escape":
			var stage: int = payload.get("stage", 0)
			print("[GameManager] Evolution fish escaped at stage %d" % stage)

## 處理幸運連鎖感染魚系統（DAY-219）
func _handle_lucky_infection_fish(payload: Dictionary) -> void:
	emit_signal("lucky_infection_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"infection_start":
			var trigger: String = payload.get("trigger_player", "")
			var total: int = payload.get("total_infected", 0)
			print("[GameManager] Infection started by %s, %d targets infected" % [trigger, total])
		"infection_spread":
			var new_count: int = payload.get("infected_targets", []).size()
			var total: int = payload.get("total_infected", 0)
			print("[GameManager] Infection spread: +%d new, total=%d" % [new_count, total])
		"infection_kill":
			var killed_id: String = payload.get("killed_target", "")
			print("[GameManager] Infected target killed: %s" % killed_id)
		"infection_blast":
			var killed: int = payload.get("total_killed", 0)
			var reward: int = payload.get("total_reward", 0)
			print("[GameManager] Infection blast: killed=%d reward=%d" % [killed, reward])

## 處理幸運反彈魚系統（DAY-220）
func _handle_lucky_ricochet_fish(payload: Dictionary) -> void:
	emit_signal("lucky_ricochet_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"ricochet_start":
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Ricochet mode started by %s" % player_name)
		"ricochet_bounce":
			var bounce_num: int = payload.get("bounce_num", 1)
			var killed: bool = payload.get("killed", false)
			print("[GameManager] Ricochet bounce #%d: killed=%s" % [bounce_num, killed])
		"ricochet_end":
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Ricochet mode ended for %s" % player_name)

## 處理幸運黑洞魚系統（DAY-221）
func _handle_lucky_black_hole(payload: Dictionary) -> void:
	emit_signal("lucky_black_hole", payload)
	var event: String = payload.get("event", "")
	match event:
		"blackhole_start":
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Black hole created by %s" % player_name)
		"blackhole_pulse":
			var pulse_num: int = payload.get("pulse_num", 1)
			var affected: int = payload.get("affected_count", 0)
			print("[GameManager] Black hole pulse #%d: affected=%d" % [pulse_num, affected])
		"singularity_blast":
			print("[GameManager] Singularity blast triggered!")
		"singularity_result":
			var killed_count: int = payload.get("killed_count", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Singularity result: killed=%d reward=%d" % [killed_count, total_reward])

## 處理幸運共鳴魚系統（DAY-222）
func _handle_lucky_resonance_fish(payload: Dictionary) -> void:
	emit_signal("lucky_resonance_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"resonance_start":
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Resonance mode started by %s" % player_name)
		"resonance_progress":
			var count: int = payload.get("count", 0)
			var target: int = payload.get("target", 30)
			print("[GameManager] Resonance progress: %d/%d" % [count, target])
		"resonance_burst":
			var total_shots: int = payload.get("total_shots", 0)
			print("[GameManager] Resonance burst! total_shots=%d" % total_shots)
		"resonance_small_burst":
			print("[GameManager] Small resonance burst triggered")
		"resonance_result":
			var affected: int = payload.get("affected_count", 0)
			var pool: int = payload.get("reward_pool", 0)
			print("[GameManager] Resonance result: affected=%d pool=%d" % [affected, pool])

## 處理幸運傳送魚系統（DAY-223）
func _handle_lucky_teleport_fish(payload: Dictionary) -> void:
	emit_signal("lucky_teleport_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"teleport_start":
			var player_name: String = payload.get("player_name", "")
			var max_waves: int = payload.get("max_waves", 4)
			print("[GameManager] Teleport vortex started by %s, max_waves=%d" % [player_name, max_waves])
		"teleport_wave":
			var wave: int = payload.get("wave", 1)
			var targets: Array = payload.get("targets", [])
			print("[GameManager] Teleport wave=%d targets=%d" % [wave, targets.size()])
			# 更新目標物位置（Client 端同步）
			_sync_teleport_positions(targets)
		"teleport_end":
			print("[GameManager] Teleport vortex ended")

## 同步傳送後的目標物位置（DAY-223）
func _sync_teleport_positions(targets: Array) -> void:
	for target_info in targets:
		var target_id: String = target_info.get("target_id", "")
		var new_x: float = target_info.get("new_x", 0.0)
		var new_y: float = target_info.get("new_y", 0.0)
		if target_id.is_empty():
			continue
		# 通知 TargetManager 更新目標位置
		emit_signal("target_teleported", target_id, Vector2(new_x, new_y))

## 處理幸運分裂魚系統（DAY-224）
func _handle_lucky_split_fish(payload: Dictionary) -> void:
	emit_signal("lucky_split_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"split_start":
			var player_name: String = payload.get("player_name", "")
			var frags: Array = payload.get("fragments", [])
			print("[GameManager] Split fish triggered by %s, fragments=%d" % [player_name, frags.size()])
		"split_blast":
			var blast_count: int = payload.get("blast_count", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Split blast: count=%d reward=%d" % [blast_count, total_reward])
		"split_end":
			print("[GameManager] Split fish ended (all fragments killed by players)")

## 處理幸運充能魚系統（DAY-225）
func _handle_lucky_charge_fish(payload: Dictionary) -> void:
	emit_signal("lucky_charge_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"charge_start":
			var target: int = payload.get("target", 10)
			print("[GameManager] Charge mode started, target=%d" % target)
		"charge_progress":
			var count: int = payload.get("count", 0)
			var target: int = payload.get("target", 10)
			print("[GameManager] Charge progress: %d/%d" % [count, target])
		"charge_ready":
			var burst_mult: float = payload.get("burst_mult", 5.0)
			print("[GameManager] Charge burst ready! mult=×%.1f" % burst_mult)
		"charge_burst":
			var reward: int = payload.get("reward", 0)
			print("[GameManager] Charge burst triggered! reward=%d" % reward)
		"charge_end":
			print("[GameManager] Charge mode ended")

## 處理幸運鏈鎖爆炸魚系統（DAY-226）
func _handle_lucky_chain_bomb(payload: Dictionary) -> void:
	emit_signal("lucky_chain_bomb", payload)
	var event: String = payload.get("event", "")
	match event:
		"chain_bomb_start":
			var marked_count: int = payload.get("marked", []).size()
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Chain bomb started by %s, marked=%d" % [player_name, marked_count])
		"chain_bomb_trigger":
			var chain_layer: int = payload.get("chain_layer", 1)
			print("[GameManager] Chain bomb triggered! layer=%d" % chain_layer)
		"chain_bomb_blast":
			var total_reward: int = payload.get("total_reward", 0)
			var chain_layer: int = payload.get("chain_layer", 1)
			print("[GameManager] Chain bomb blast! layer=%d reward=%d" % [chain_layer, total_reward])
		"chain_bomb_expire":
			print("[GameManager] Chain bomb marks expired")

## 處理幸運鏡像時空魚系統（DAY-227）
func _handle_lucky_mirror_time(payload: Dictionary) -> void:
	emit_signal("lucky_mirror_time", payload)
	var event: String = payload.get("event", "")
	match event:
		"time_rewind_start":
			var rewind_count: int = payload.get("rewind_count", 0)
			var boost_mult: float = payload.get("boost_mult", 2.0)
			print("[GameManager] Time rewind started! rewound=%d boost=×%.1f" % [rewind_count, boost_mult])
		"time_collapse":
			var collapse_count: int = payload.get("collapse_count", 0)
			print("[GameManager] Time collapse! %d targets HP -40%%" % collapse_count)

## 處理幸運量子魚系統（DAY-228）
func _handle_lucky_quantum_fish(payload: Dictionary) -> void:
	emit_signal("lucky_quantum_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"quantum_start":
			var quantum_count: int = payload.get("quantum_count", 0)
			var high_mult: float = payload.get("high_mult", 3.0)
			var low_mult: float = payload.get("low_mult", 0.8)
			print("[GameManager] Quantum superposition! count=%d high=×%.1f low=×%.1f" % [quantum_count, high_mult, low_mult])
		"quantum_collapse":
			var collapse_high: bool = payload.get("collapse_high", false)
			var collapse_mult: float = payload.get("collapse_mult", 1.0)
			print("[GameManager] Quantum collapse! high=%s mult=×%.1f" % [str(collapse_high), collapse_mult])
		"quantum_blast":
			var blast_count: int = payload.get("blast_count", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Quantum blast! count=%d reward=%d" % [blast_count, total_reward])

## 處理幸運寄生魚系統（DAY-229）
func _handle_lucky_parasite_fish(payload: Dictionary) -> void:
	emit_signal("lucky_parasite_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"parasite_start":
			var parasite_count: int = payload.get("parasite_count", 0)
			var kill_mult: float = payload.get("kill_mult", 2.2)
			print("[GameManager] Parasite released! count=%d mult=×%.1f" % [parasite_count, kill_mult])
		"parasite_tick":
			var target_id: String = payload.get("target_id", "")
			var tick_count: int = payload.get("tick_count", 1)
			print("[GameManager] Parasite tick=%d target=%s" % [tick_count, target_id])
		"parasite_jump":
			var jump_layer: int = payload.get("jump_layer", 1)
			print("[GameManager] Parasite jumped! layer=%d" % jump_layer)
		"parasite_kill":
			var kill_reward: int = payload.get("kill_reward", 0)
			print("[GameManager] Parasite target killed! reward=%d" % kill_reward)
		"parasite_end":
			var target_id: String = payload.get("target_id", "")
			print("[GameManager] Parasite expired on target=%s" % target_id)

## 處理幸運風暴魚系統（DAY-230）
func _handle_lucky_storm_fish(payload: Dictionary) -> void:
	emit_signal("lucky_storm_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"storm_start":
			var storm_x: float = payload.get("storm_x", 500.0)
			var storm_y: float = payload.get("storm_y", 300.0)
			var kill_mult: float = payload.get("kill_mult", 2.5)
			print("[GameManager] Storm started at (%.0f,%.0f) mult=×%.1f" % [storm_x, storm_y, kill_mult])
		"storm_rotate":
			var rotate_count: int = payload.get("rotate_count", 1)
			var moved_count: int = payload.get("moved_count", 0)
			print("[GameManager] Storm rotate #%d moved=%d" % [rotate_count, moved_count])
		"storm_blast_start":
			print("[GameManager] Storm blast starting!")
		"storm_blast":
			var killed_count: int = payload.get("killed_count", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Storm blast! killed=%d reward=%d" % [killed_count, total_reward])

## 處理幸運迴旋鏢魚系統（DAY-231）
func _handle_lucky_boomerang_fish(payload: Dictionary) -> void:
	emit_signal("lucky_boomerang_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"boomerang_start":
			var player_name: String = payload.get("player_name", "")
			var max_bounces: int = payload.get("max_bounces", 3)
			print("[GameManager] Boomerang mode started! player=%s max_bounces=%d" % [player_name, max_bounces])
		"boomerang_hit":
			var bounce_num: int = payload.get("bounce_num", 1)
			var killed: bool = payload.get("killed", false)
			var reward: int = payload.get("reward", 0)
			print("[GameManager] Boomerang hit! bounce#%d killed=%s reward=%d" % [bounce_num, str(killed), reward])
		"boomerang_end":
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Boomerang mode ended! player=%s" % player_name)

func _handle_lucky_magnet_fish(payload: Dictionary) -> void:
	emit_signal("lucky_magnet_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"magnet_start":
			var player_name: String = payload.get("player_name", "")
			var kill_boost: float = payload.get("kill_boost", 1.8)
			print("[GameManager] Magnet field started! player=%s boost=x%.1f" % [player_name, kill_boost])
		"magnet_pull":
			var pull_num: int = payload.get("pull_num", 1)
			var moved_count: int = payload.get("moved_count", 0)
			print("[GameManager] Magnet pull#%d: %d targets moved" % [pull_num, moved_count])
			# 同步目標物位置（磁力吸引後更新 Client 端位置）
			var positions: Array = payload.get("positions", [])
			_sync_magnet_positions(positions)
		"magnet_blast":
			var killed_count: int = payload.get("killed_count", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Magnet blast! killed=%d total_reward=%d" % [killed_count, total_reward])
		"magnet_end":
			print("[GameManager] Magnet field ended!")

## 同步磁力吸引後的目標物位置（DAY-232）
func _sync_magnet_positions(positions: Array) -> void:
	for pos_info in positions:
		var target_id: String = pos_info.get("id", "")
		var new_x: float = pos_info.get("x", 0.0)
		var new_y: float = pos_info.get("y", 0.0)
		if target_id.is_empty():
			continue
		# 複用傳送魚的 target_teleported 訊號同步位置（平滑移動）
		emit_signal("target_teleported", target_id, Vector2(new_x, new_y))

## 幸運回聲魚系統（DAY-233）
func _handle_lucky_echo_fish(payload: Dictionary) -> void:
	emit_signal("lucky_echo_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"echo_ready":
			print("[GameManager] Echo mode activated! (personal)")
		"echo_broadcast":
			var player_name: String = payload.get("player_name", "")
			print("[GameManager] Echo mode broadcast! player=%s" % player_name)
		"echo_spawn":
			var layer_num: int = payload.get("layer", 1)
			var echo_instance_id: String = payload.get("echo_instance_id", "")
			var mult_label: String = payload.get("mult_label", "×1.5")
			print("[GameManager] Echo spawn! layer=%d instanceID=%s mult=%s" % [layer_num, echo_instance_id, mult_label])
		"echo_spawn_personal":
			var layer_num: int = payload.get("layer", 1)
			var mult_label: String = payload.get("mult_label", "×1.5")
			print("[GameManager] Echo spawn personal! layer=%d mult=%s" % [layer_num, mult_label])
		"echo_expire":
			var echo_instance_id: String = payload.get("echo_instance_id", "")
			print("[GameManager] Echo expired! instanceID=%s" % echo_instance_id)

## 幸運漩渦魚系統（DAY-234）
func _handle_lucky_vortex_fish(payload: Dictionary) -> void:
	emit_signal("lucky_vortex_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"vortex_start":
			var player_name: String = payload.get("player_name", "")
			var kill_boost: float = payload.get("kill_boost", 2.2)
			print("[GameManager] Vortex started! player=%s boost=x%.1f" % [player_name, kill_boost])
		"vortex_rotate":
			var rotate_num: int = payload.get("rotate_num", 1)
			var rotated_count: int = payload.get("rotated_count", 0)
			print("[GameManager] Vortex rotate#%d: %d targets rotated" % [rotate_num, rotated_count])
			# 同步目標物旋轉後的位置
			var positions: Array = payload.get("positions", [])
			_sync_vortex_positions(positions)
		"vortex_blast":
			var killed_count: int = payload.get("killed_count", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Vortex blast! killed=%d total_reward=%d" % [killed_count, total_reward])
		"vortex_end":
			print("[GameManager] Vortex ended!")

## 同步漩渦旋轉後的目標物位置（DAY-234）
func _sync_vortex_positions(positions: Array) -> void:
	for pos_info in positions:
		var target_id: String = pos_info.get("id", "")
		var new_x: float = pos_info.get("x", 0.0)
		var new_y: float = pos_info.get("y", 0.0)
		if target_id.is_empty():
			continue
		# 複用 target_teleported 訊號同步位置（平滑移動）
		emit_signal("target_teleported", target_id, Vector2(new_x, new_y))

## 幸運時間炸彈魚系統（DAY-235）
func _handle_lucky_time_bomb_fish(payload: Dictionary) -> void:
	emit_signal("lucky_time_bomb_fish", payload)
	var event: String = payload.get("event", "")
	match event:
		"bomb_placed":
			var player_name: String = payload.get("player_name", "")
			var bomb_count: int = payload.get("bomb_count", 4)
			var fuse_sec: int = payload.get("fuse_sec", 8)
			print("[GameManager] Time bombs placed! player=%s count=%d fuse=%ds" % [player_name, bomb_count, fuse_sec])
		"bomb_countdown":
			var target_id: String = payload.get("target_id", "")
			var remaining: int = payload.get("remaining", 0)
			print("[GameManager] Bomb countdown: target=%s remaining=%d" % [target_id, remaining])
		"bomb_early_detonate":
			var player_name: String = payload.get("player_name", "")
			var mult: float = payload.get("mult", 2.0)
			print("[GameManager] Early detonation! player=%s mult=x%.1f" % [player_name, mult])
		"bomb_chain_blast":
			var killed_count: int = payload.get("killed_count", 0)
			var total_reward: int = payload.get("total_reward", 0)
			print("[GameManager] Chain blast! killed=%d reward=%d" % [killed_count, total_reward])
		"bomb_auto_explode":
			var target_id: String = payload.get("target_id", "")
			var killed: bool = payload.get("killed", false)
			var reward: int = payload.get("reward", 0)
			print("[GameManager] Auto explode: target=%s killed=%s reward=%d" % [target_id, str(killed), reward])

## 幸運鏡面世界魚系統（DAY-236）
func _handle_lucky_mirror_world(payload: Dictionary) -> void:
	emit_signal("lucky_mirror_world", payload)
	var event: String = payload.get("event", "")
	match event:
		"mirror_start":
			var player_name: String = payload.get("player_name", "")
			var kill_boost: float = payload.get("kill_boost", 2.3)
			print("[GameManager] Mirror world started! player=%s boost=x%.1f" % [player_name, kill_boost])
			# 同步所有目標物的鏡像後位置
			var positions: Array = payload.get("positions", [])
			_sync_mirror_positions(positions)
		"mirror_collapse":
			var collapsed_count: int = payload.get("collapsed_count", 0)
			print("[GameManager] Mirror collapse! affected=%d targets" % collapsed_count)
		"mirror_end":
			print("[GameManager] Mirror world ended!")

## 同步鏡像反轉後的目標物位置（DAY-236）
func _sync_mirror_positions(positions: Array) -> void:
	for pos_info in positions:
		var target_id: String = pos_info.get("id", "")
		var new_x: float = pos_info.get("x", 0.0)
		var new_y: float = pos_info.get("y", 0.0)
		if target_id.is_empty():
			continue
		# 複用 target_teleported 訊號同步位置（平滑移動）
		emit_signal("target_teleported", target_id, Vector2(new_x, new_y))

## 幸運冰凍世界魚系統（DAY-237）
func _handle_lucky_freeze_world(payload: Dictionary) -> void:
	emit_signal("lucky_freeze_world", payload)
	var event: String = payload.get("event", "")
	match event:
		"freeze_start":
			var player_name: String = payload.get("player_name", "")
			var kill_boost: float = payload.get("kill_boost", 2.0)
			var frozen_count: int = payload.get("frozen_count", 0)
			var speed_factor: float = payload.get("speed_factor", 0.2)
			print("[GameManager] Freeze world started! player=%s boost=x%.1f frozen=%d speed_factor=%.1f" % [player_name, kill_boost, frozen_count, speed_factor])
			# 通知 TargetManager 降低所有目標速度
			emit_signal("freeze_world_started", speed_factor)
		"freeze_crack":
			var cracked_count: int = payload.get("cracked_count", 0)
			print("[GameManager] Ice crack! affected=%d targets" % cracked_count)
		"freeze_end":
			print("[GameManager] Freeze world ended! Speed restored.")
			# 通知 TargetManager 恢復所有目標速度
			emit_signal("freeze_world_ended")

## 幸運重力反轉魚系統（DAY-238）
func _handle_lucky_gravity_flip(payload: Dictionary) -> void:
	emit_signal("lucky_gravity_flip", payload)
	var event: String = payload.get("event", "")
	match event:
		"gravity_start":
			var player_name: String = payload.get("player_name", "")
			var kill_boost: float = payload.get("kill_boost", 2.1)
			var positions = payload.get("positions", [])
			print("[GameManager] Gravity flip started! player=%s boost=x%.1f targets=%d" % [player_name, kill_boost, positions.size()])
			# 同步目標 Y 座標翻轉位置
			_sync_gravity_positions(positions)
		"gravity_collapse":
			var collapsed_count: int = payload.get("collapsed_count", 0)
			print("[GameManager] Gravity collapse! affected=%d targets" % collapsed_count)
		"gravity_end":
			print("[GameManager] Gravity flip ended! Positions restored.")

## 同步重力翻轉後的目標位置（DAY-238）
func _sync_gravity_positions(positions: Array) -> void:
	for pos_data in positions:
		var target_id: String = pos_data.get("id", "")
		var new_x: float = pos_data.get("x", 0.0)
		var new_y: float = pos_data.get("y", 0.0)
		if target_id != "":
			emit_signal("target_teleported", target_id, new_x, new_y)
