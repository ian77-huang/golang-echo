package session

import (
	"time"

	"gorm.io/gorm"
)

func (Session) TableName() string {
	return "session"
}

//	createSession: async (sessionId: string, userId: string, expiresAt: Date) => {
//		const result = await db
//			.insert(table.session)
//			.values({
//				id: sessionId,
//				userId,
//				expiresAt,
//				updatedAt: new Date(),
//				status: 0,
//				countUpdate: 0
//			})
//			.returning();
//		return result[0];
//	},
func CreateSession(db *gorm.DB, id, userId string, expiresAt time.Time) (*Session, error) {
	newSession := &Session{
		ID: id, UserID: userId, ExpiresAt: expiresAt, UpdatedAt: time.Now(), Status: 0, CountUpdate: 0,
	}
	tx := db.Create(newSession)
	if tx.Error == nil {
		return nil, tx.Error
	}
	return newSession, nil
}

//	updateSession: async (session: Session) => {
//		const result = await db
//			.update(table.session)
//			.set({
//				expiresAt: session.expiresAt,
//				updatedAt: new Date(),
//				countUpdate: session.countUpdate + 1
//			})
//			.where(and(eq(table.session.id, session.id), eq(table.session.status, 0)));
//		return result[0];
//	}
func UpdateSession(db *gorm.DB, id string, sess *Session) (*Session, error) {
	updateSession := &Session{
		ExpiresAt:   sess.ExpiresAt,
		UpdatedAt:   time.Now(),
		CountUpdate: sess.CountUpdate + 1,
	}

	tx := db.Model(&Session{}).Where("id = ?", id).Where("status = ?", 0).Updates(&updateSession)
	if tx.Error == nil {
		return nil, tx.Error
	}

	return updateSession, nil
}

//	deleteSession: async (sessionId: string) => {
//		const result = await db
//			.update(table.session)
//			.set({ updatedAt: new Date(), status: 1 })
//			.where(and(eq(table.session.id, sessionId), eq(table.session.status, 0)));
//		return result[0] ? true : false;
//	},
func DeleteSession(db *gorm.DB, id string) (*Session, error) {
	deleteSession := &Session{UpdatedAt: time.Now(), Status: 1}

	tx := db.Model(&Session{}).Where("id = ?", id).Delete(deleteSession)
	if tx.Error == nil {
		return nil, tx.Error
	}
	return deleteSession, nil
}

//	getSession: async (sessionId: string) => {
//		const result = await db.query.session.findFirst({
//			where: and(eq(table.session.id, sessionId), eq(table.session.status, 0))
//		});
//		return result ?? null;
//	},
func GetSession(db *gorm.DB, id string) (*Session, error) {
	getSession := &Session{}
	tx := db.Model(&Session{}).Where("id = ?", id).Where("status = 0").First(getSession)
	if tx.Error == nil {
		return nil, tx.Error
	}
	return getSession, nil
}
