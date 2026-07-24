package domain

// IsAssignableStaffRole reports whether a role may be assigned to clinic staff accounts.
func IsAssignableStaffRole(role Role) bool {
	switch role {
	case RoleAdmin, RoleDoctor, RoleNurse, RoleStaff:
		return true
	default:
		return false
	}
}

// CanReadPatientChart reports whether a role may view patient records.
func CanReadPatientChart(role Role) bool {
	switch role {
	case RoleAdmin, RoleDoctor, RoleNurse, RoleStaff:
		return true
	default:
		return false
	}
}

// CanEditPatientDemographics reports whether a role may create or update patient profiles.
func CanEditPatientDemographics(role Role) bool {
	switch role {
	case RoleAdmin, RoleDoctor, RoleNurse:
		return true
	default:
		return false
	}
}

// CanDeletePatient reports whether a role may delete a patient record.
func CanDeletePatient(role Role) bool {
	return role == RoleAdmin
}

// CanManageClinicalNotes reports whether a role may create, edit, or delete clinical notes.
func CanManageClinicalNotes(role Role) bool {
	switch role {
	case RoleAdmin, RoleDoctor, RoleNurse:
		return true
	default:
		return false
	}
}

// CanUploadDocuments reports whether a role may upload patient documents.
func CanUploadDocuments(role Role) bool {
	switch role {
	case RoleAdmin, RoleDoctor, RoleNurse, RoleStaff:
		return true
	default:
		return false
	}
}

// CanDeleteDocuments reports whether a role may delete patient documents.
func CanDeleteDocuments(role Role) bool {
	switch role {
	case RoleAdmin, RoleDoctor:
		return true
	default:
		return false
	}
}

// CanSubmitClinicalForms reports whether a role may submit structured clinical forms.
func CanSubmitClinicalForms(role Role) bool {
	switch role {
	case RoleAdmin, RoleDoctor, RoleNurse:
		return true
	default:
		return false
	}
}

// ChartReaderRoles returns roles allowed to read patient charts.
func ChartReaderRoles() []Role {
	return []Role{RoleAdmin, RoleDoctor, RoleNurse, RoleStaff}
}

// ClinicalStaffRoles returns roles allowed to perform clinical chart edits.
func ClinicalStaffRoles() []Role {
	return []Role{RoleAdmin, RoleDoctor, RoleNurse}
}

// DocumentUploaderRoles returns roles allowed to upload documents.
func DocumentUploaderRoles() []Role {
	return []Role{RoleAdmin, RoleDoctor, RoleNurse, RoleStaff}
}

// DocumentDeleterRoles returns roles allowed to delete documents.
func DocumentDeleterRoles() []Role {
	return []Role{RoleAdmin, RoleDoctor}
}
