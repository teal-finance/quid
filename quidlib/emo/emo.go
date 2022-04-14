package emo

// Info :.
func (zone Zone) Info(args ...interface{}) Event {
	return processEvent("ℹ️", zone, false, args)
}

// Warning :.
func (zone Zone) Warning(args ...interface{}) Event {
	return processEvent("🔔", zone, true, args)
}

// Error :.
func (zone Zone) Error(args ...interface{}) Event {
	return processEvent("💢", zone, true, args)
}

// Query :.
func (zone Zone) Query(args ...interface{}) Event {
	return processEvent("🗄️", zone, false, args)
}

// QueryError :.
func (zone Zone) QueryError(args ...interface{}) Event {
	return processEvent("🗄️", zone, true, args)
}

// Encrypt :.
func (zone Zone) Encrypt(args ...interface{}) Event {
	return processEvent("🎼", zone, false, args)
}

// EncryptError :.
func (zone Zone) EncryptError(args ...interface{}) Event {
	return processEvent("🎼", zone, true, args)
}

// Decrypt :.
func (zone Zone) Decrypt(args ...interface{}) Event {
	return processEvent("🗝️", zone, false, args)
}

// DecryptError :.
func (zone Zone) DecryptError(args ...interface{}) Event {
	return processEvent("🗝️", zone, true, args)
}

// Time :.
func (zone Zone) Time(args ...interface{}) Event {
	return processEvent("⏱️", zone, false, args)
}

// TimeError :.
func (zone Zone) TimeError(args ...interface{}) Event {
	return processEvent("⏱️", zone, true, args)
}

// Param :.
func (zone Zone) Param(args ...interface{}) Event {
	return processEvent("📥", zone, false, args)
}

// ParamError :.
func (zone Zone) ParamError(args ...interface{}) Event {
	return processEvent("📥", zone, true, args)
}

// Debug :.
func (zone Zone) Debug(args ...interface{}) Event {
	return processEvent("💊", zone, false, args)
}
