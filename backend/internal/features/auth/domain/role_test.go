package domain

import "testing"

func TestIsAssignableStaffRole(t *testing.T) {
	assignable := []Role{RoleAdmin, RoleDoctor, RoleNurse, RoleStaff}
	for _, role := range assignable {
		if !IsAssignableStaffRole(role) {
			t.Fatalf("expected %q to be assignable", role)
		}
	}

	if IsAssignableStaffRole(RolePatient) {
		t.Fatal("patient role must not be assignable to staff accounts")
	}
	if IsAssignableStaffRole(Role("superuser")) {
		t.Fatal("unknown role must not be assignable")
	}
}

func TestPatientAccessRoles(t *testing.T) {
	if !CanDeletePatient(RoleAdmin) {
		t.Fatal("admin should delete patients")
	}
	if CanDeletePatient(RoleStaff) {
		t.Fatal("staff should not delete patients")
	}
	if !CanReadPatientChart(RoleStaff) {
		t.Fatal("staff should read charts")
	}
	if CanEditPatientDemographics(RoleStaff) {
		t.Fatal("staff should not edit demographics")
	}
}

func TestClinicalChartWriteRoles(t *testing.T) {
	// Exercise log writes use ClinicalStaffRoles (same gate as clinical notes/forms).
	for _, role := range []Role{RoleAdmin, RoleDoctor, RoleNurse} {
		if !CanManageClinicalNotes(role) || !CanSubmitClinicalForms(role) {
			t.Fatalf("%q should manage clinical chart writes", role)
		}
	}
	if CanManageClinicalNotes(RoleStaff) || CanSubmitClinicalForms(RoleStaff) {
		t.Fatal("front-desk staff must not write clinical notes/forms/exercise log")
	}

	clinical := ClinicalStaffRoles()
	if len(clinical) != 3 {
		t.Fatalf("expected 3 clinical roles, got %d", len(clinical))
	}
	readers := ChartReaderRoles()
	if len(readers) != 4 {
		t.Fatalf("expected 4 chart reader roles, got %d", len(readers))
	}
}
