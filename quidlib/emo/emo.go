package emo

// Query :
func (zone Zone) Query(errObjs ...interface{}) Event {
	return processEvent("🗄️", zone, false, errObjs)
}

// QueryError :
func (zone Zone) QueryError(errObjs ...interface{}) Event {
	return processEvent("🗄️", zone, true, errObjs)
}

// Encrypt :
func (zone Zone) Encrypt(errObjs ...interface{}) Event {
	return processEvent("🎼", zone, false, errObjs)
}

// EncryptError :
func (zone Zone) EncryptError(errObjs ...interface{}) Event {
	return processEvent("🎼", zone, true, errObjs)
}

// Decrypt :
func (zone Zone) Decrypt(errObjs ...interface{}) Event {
	return processEvent("🗝️", zone, false, errObjs)
}

// DecryptError :
func (zone Zone) DecryptError(errObjs ...interface{}) Event {
	return processEvent("🗝️", zone, true, errObjs)
}

// Time :
func (zone Zone) Time(errObjs ...interface{}) Event {
	return processEvent("⏱️", zone, false, errObjs)
}

// TimeError :
func (zone Zone) TimeError(errObjs ...interface{}) Event {
	return processEvent("⏱️", zone, true, errObjs)
}

// Param :
func (zone Zone) Param(errObjs ...interface{}) Event {
	return processEvent("📥", zone, false, errObjs)
}

// ParamError :
func (zone Zone) ParamError(errObjs ...interface{}) Event {
	return processEvent("📥", zone, true, errObjs)
}

// Debug :
func (zone Zone) Debug(errObjs ...interface{}) Event {
	return processEvent("💊", zone, false, errObjs)
}
