package tests

import (
	"testing"

	"github.com/matrix-org/complement/internal/b"
)

// TestKnockingLocal tests that a user knocking on a room which the homeserver is already a part of works
func TestKnockingLocal(t *testing.T) {
	deployment := Deploy(t, "local_knocking", b.BlueprintAliceBob)
	defer deployment.Destroy(t)

	aliceUserID := "@alice:hs1"
	alice := deployment.Client(t, "hs1", aliceUserID)
	roomID := alice.CreateRoom(t, struct {
		Preset      string `json:"preset"`
		RoomVersion string `json:"room_version"`
		Name        string `json:"name"`
		Topic       string `json:"topic"`
	}{
		"private",           // Set to private in order to get an invite-only room
		"xyz.amorgan.knock", // Room version required for knocking. TODO: Remove when knocking is in a stable room version
		// Add some state to the room. We'll later check that this comes down sync correctly.
		"knocking test room",
		"Who's there?",
	})

	//bobUserID := "@bob:hs1"
	//bob := deployment.Client(t, "hs1", bobUserID)

	t.Run("Set the join_rules of a private room to 'knock'", func(t *testing.T) {
		alice.MustDo(
			t,
			"PUT",
			[]string{"_matrix", "client", "r0", "rooms", roomID, "state", "m.room.join_rules", ""},
			struct {
				JoinRule string `json:"join_rule"`
			}{
				"knock",
			},
		)
	})

	/*
		t.Run("parallel", func(t *testing.T) {
			t.Run("", func(t *testing.T) {
				t.Parallel()
				alice := deployment.Client(t, "hs1", userID)
				alice.SyncUntilTimelineHas(t, roomID, func(ev gjson.Result) bool {
					if ev.Get("type").Str != "m.room.create" {
						return false
					}
					must.EqualStr(t, ev.Get("sender").Str, userID, "wrong sender")
					must.EqualStr(t, ev.Get("content").Get("creator").Str, userID, "wrong content.creator")
					return true
				})
			})
			// sytest: Room creation reports m.room.member to myself
			t.Run("Room creation reports m.room.member to myself", func(t *testing.T) {
				t.Parallel()
				alice := deployment.Client(t, "hs1", userID)
				alice.SyncUntilTimelineHas(t, roomID, func(ev gjson.Result) bool {
					if ev.Get("type").Str != "m.room.member" {
						return false
					}
					must.EqualStr(t, ev.Get("sender").Str, userID, "wrong sender")
					must.EqualStr(t, ev.Get("state_key").Str, userID, "wrong state_key")
					must.EqualStr(t, ev.Get("content").Get("membership").Str, "join", "wrong content.membership")
					return true
				})
			})
		})
	*/
}
