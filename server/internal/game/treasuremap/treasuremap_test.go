package treasuremap

import (
	"testing"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
}

func TestRecordKill_UnknownTarget(t *testing.T) {
	m := New()
	filled, lines, full := m.RecordKill("p1", "UNKNOWN")
	if filled || len(lines) > 0 || full {
		t.Error("unknown target should not fill any cell")
	}
}

func TestRecordKill_FillCell(t *testing.T) {
	m := New()
	filled, lines, full := m.RecordKill("p1", "T001") // row=0, col=0
	if !filled {
		t.Error("T001 should fill cell (0,0)")
	}
	if len(lines) > 0 {
		t.Error("single cell should not complete a line")
	}
	if full {
		t.Error("single cell should not complete the map")
	}
}

func TestRecordKill_NoDuplicate(t *testing.T) {
	m := New()
	m.RecordKill("p1", "T001")
	filled, _, _ := m.RecordKill("p1", "T001") // 重複擊破
	if filled {
		t.Error("duplicate kill should not fill cell again")
	}
}

func TestRecordKill_CompleteLine_Row(t *testing.T) {
	m := New()
	// 填滿第一行：T001(0,0), T003(0,1), T005(0,2)
	m.RecordKill("p1", "T001")
	m.RecordKill("p1", "T003")
	_, lines, _ := m.RecordKill("p1", "T005")
	if len(lines) == 0 {
		t.Error("should complete row0")
	}
	if lines[0].Type != "row0" {
		t.Errorf("expected row0, got %s", lines[0].Type)
	}
}

func TestRecordKill_CompleteLine_Col(t *testing.T) {
	m := New()
	// 填滿第一列：T001(0,0), T002(1,0), T006(2,0)
	m.RecordKill("p1", "T001")
	m.RecordKill("p1", "T002")
	_, lines, _ := m.RecordKill("p1", "T006")
	if len(lines) == 0 {
		t.Error("should complete col0")
	}
	if lines[0].Type != "col0" {
		t.Errorf("expected col0, got %s", lines[0].Type)
	}
}

func TestRecordKill_CompleteLine_Diag(t *testing.T) {
	m := New()
	// 填滿對角線：T001(0,0), T101(1,1), T105(2,2)
	m.RecordKill("p1", "T001")
	m.RecordKill("p1", "T101")
	_, lines, _ := m.RecordKill("p1", "T105")
	if len(lines) == 0 {
		t.Error("should complete diag0")
	}
	if lines[0].Type != "diag0" {
		t.Errorf("expected diag0, got %s", lines[0].Type)
	}
}

func TestRecordKill_FullMap(t *testing.T) {
	m := New()
	// 填滿所有 9 格
	allTargets := []string{"T001", "T003", "T005", "T002", "T101", "T104", "T006", "T102", "T105"}
	var fullDone bool
	for _, id := range allTargets {
		_, _, fullDone = m.RecordKill("p1", id)
	}
	if !fullDone {
		t.Error("should complete full map after filling all 9 cells")
	}
}

func TestRecordKill_FullMap_NoDuplicate(t *testing.T) {
	m := New()
	allTargets := []string{"T001", "T003", "T005", "T002", "T101", "T104", "T006", "T102", "T105"}
	for _, id := range allTargets {
		m.RecordKill("p1", id)
	}
	// 再次填滿不應再觸發 fullDone
	_, _, fullDone := m.RecordKill("p1", "T001")
	if fullDone {
		t.Error("full map should not trigger again")
	}
}

func TestGetSnapshot(t *testing.T) {
	m := New()
	m.RecordKill("p1", "T001")
	snap := m.GetSnapshot("p1")
	if snap == nil {
		t.Fatal("snapshot should not be nil")
	}
	if !snap.Cells[0][0] {
		t.Error("cell (0,0) should be filled")
	}
	if snap.FilledCount != 1 {
		t.Errorf("filled count should be 1, got %d", snap.FilledCount)
	}
}

func TestRemovePlayer(t *testing.T) {
	m := New()
	m.RecordKill("p1", "T001")
	m.RemovePlayer("p1")
	snap := m.GetSnapshot("p1")
	if snap.FilledCount != 0 {
		t.Error("player data should be cleared after remove")
	}
}

func TestCalcLineReward(t *testing.T) {
	reward := CalcLineReward(10) // LV5 betCost=10
	if reward != 500 {
		t.Errorf("line reward should be 500, got %d", reward)
	}
}

func TestCalcFullReward(t *testing.T) {
	reward := CalcFullReward(10) // LV5 betCost=10
	if reward != 5000 {
		t.Errorf("full reward should be 5000, got %d", reward)
	}
}

func TestGetCellDef(t *testing.T) {
	def := GetCellDef(0, 0)
	if def == nil {
		t.Fatal("cell def should not be nil")
	}
	if def.DefID != "T001" {
		t.Errorf("cell (0,0) should be T001, got %s", def.DefID)
	}
}

func TestGetCellDef_OutOfBounds(t *testing.T) {
	def := GetCellDef(-1, 0)
	if def != nil {
		t.Error("out of bounds should return nil")
	}
	def = GetCellDef(3, 0)
	if def != nil {
		t.Error("out of bounds should return nil")
	}
}
