package msg

import (
	"fmt"
	"github.com/Blackrush/gofus/realm/db"
	sdb "github.com/Blackrush/gofus/shared/db"
	"io"
	"strings"
	"time"
)

type HelloGame struct{}

func (msg *HelloGame) Opcode() string                 { return "HG" }
func (msg *HelloGame) Serialize(out io.Writer) error  { return nil }
func (msg *HelloGame) Deserialize(in io.Reader) error { return nil }

type RealmLoginReq struct {
	Ticket string
}

func (msg *RealmLoginReq) Opcode() string { return "AT" }
func (msg *RealmLoginReq) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Ticket)
	return nil
}
func (msg *RealmLoginReq) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "AT%s", &msg.Ticket)
	return nil
}

type RealmLoginSuccess struct {
	CommunityId int
}

func (msg *RealmLoginSuccess) Opcode() string { return "ATK" }
func (msg *RealmLoginSuccess) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.CommunityId)
	return nil
}
func (msg *RealmLoginSuccess) Deserialize(in io.Reader) error { return nil }

type RealmLoginError struct{}

func (msg *RealmLoginError) Opcode() string                 { return "ATE" }
func (msg *RealmLoginError) Serialize(out io.Writer) error  { return nil }
func (msg *RealmLoginError) Deserialize(in io.Reader) error { return nil }

type ClientUseKeyReq struct {
	KeyId int
}

func (msg *ClientUseKeyReq) Opcode() string { return "Ak" }
func (msg *ClientUseKeyReq) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%x", msg.KeyId)
	return nil
}
func (msg *ClientUseKeyReq) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "Ak%x", &msg.KeyId)
	return nil
}

type RegionalVersionReq struct{}

func (msg *RegionalVersionReq) Opcode() string                 { return "AV" }
func (msg *RegionalVersionReq) Serialize(out io.Writer) error  { return nil }
func (msg *RegionalVersionReq) Deserialize(in io.Reader) error { return nil }

type RegionalVersionResp struct {
	CommunityId int
}

func (msg *RegionalVersionResp) Opcode() string { return "AV" }
func (msg *RegionalVersionResp) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.CommunityId)
	return nil
}
func (msg *RegionalVersionResp) Deserialize(in io.Reader) error { return nil }

type PlayersGiftsReq struct {
	Language string
}

func (msg *PlayersGiftsReq) Opcode() string { return "Ag" }
func (msg *PlayersGiftsReq) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Language)
	return nil
}
func (msg *PlayersGiftsReq) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "Ag%s", &msg.Language)
	return nil
}

type SetPlayerGiftReq struct {
	GiftId   int
	PlayerId int
}

func (msg *SetPlayerGiftReq) Opcode() string { return "AG" }
func (msg *SetPlayerGiftReq) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%d|%d", msg.GiftId, msg.PlayerId)
	return nil
}
func (msg *SetPlayerGiftReq) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "%d|%d", &msg.GiftId, &msg.PlayerId)
	return nil
}

type SetIdentity struct {
	Identity string
}

func (msg *SetIdentity) Opcode() string { return "Ai" }
func (msg *SetIdentity) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Identity)
	return nil
}
func (msg *SetIdentity) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "Ai%s", &msg.Identity)
	return nil
}

type PlayersReq struct{}

func (msg *PlayersReq) Opcode() string                 { return "AL" }
func (msg *PlayersReq) Serialize(out io.Writer) error  { return nil }
func (msg *PlayersReq) Deserialize(in io.Reader) error { return nil }

type PlayersResp struct {
	ServerId        uint
	SubscriptionEnd time.Time
	Players         []*db.Player
}

func (msg *PlayersResp) Opcode() string { return "ALK" }
func (msg *PlayersResp) Serialize(out io.Writer) error {
	var remainingSubscription int64
	diff := msg.SubscriptionEnd.Sub(time.Now())
	if diff.Nanoseconds() > 0 {
		remainingSubscription = diff.Nanoseconds() / 1e6
	}

	fmt.Fprintf(out, "%d|%d", remainingSubscription, len(msg.Players))
	for _, player := range msg.Players {
		fmt.Fprintf(out, "|%d;%s;%d;%d;%v;%v;%v;%s;%d;;;",
			player.Id,
			player.Name,
			player.Experience.Level,
			player.Appearance.Skin,
			player.Appearance.Colors.First,
			player.Appearance.Colors.Second,
			player.Appearance.Colors.Third,
			player.Appearance.Accessories,
			msg.ServerId,
		)
	}
	return nil
}
func (msg *PlayersResp) Deserialize(in io.Reader) error {
	return nil
}

type RandNameReq struct{}

func (msg *RandNameReq) Opcode() string                 { return "AP" }
func (msg *RandNameReq) Serialize(out io.Writer) error  { return nil }
func (msg *RandNameReq) Deserialize(in io.Reader) error { return nil }

type RandNameResp struct {
	Name string
}

func (msg *RandNameResp) Opcode() string { return "AP" }
func (msg *RandNameResp) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Name)
	return nil
}
func (msg *RandNameResp) Deserialize(in io.Reader) error {
	return nil
}

type CreatePlayerReq struct {
	Name   string
	Breed  int
	Gender bool
	Colors struct {
		First  int
		Second int
		Third  int
	}
}

func (msg *CreatePlayerReq) Opcode() string { return "AA" }
func (msg *CreatePlayerReq) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%s|%d|%d|%d|%d|%d", msg.Name, msg.Breed, btoi(msg.Gender), msg.Colors.First, msg.Colors.Second, msg.Colors.Third)
	return nil
}
func (msg *CreatePlayerReq) Deserialize(in io.Reader) error {
	var body string
	fmt.Fscanf(in, "AA%s", &body)

	args := strings.SplitN(body, "|", 6)

	msg.Name = args[0]
	msg.Breed = atoi(args[1])
	msg.Gender = aitob(args[2])
	msg.Colors.First = atoi(args[3])
	msg.Colors.Second = atoi(args[4])
	msg.Colors.Third = atoi(args[5])

	return nil
}

type CreatePlayerErrorResp struct{}

func (msg *CreatePlayerErrorResp) Opcode() string                 { return "AAE" }
func (msg *CreatePlayerErrorResp) Serialize(out io.Writer) error  { return nil }
func (msg *CreatePlayerErrorResp) Deserialize(in io.Reader) error { return nil }

type PlayerSelectionReq struct {
	PlayerId uint64
}

func (msg *PlayerSelectionReq) Opcode() string { return "AS" }
func (msg *PlayerSelectionReq) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%d", msg.PlayerId)
	return nil
}
func (msg *PlayerSelectionReq) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "AS%d", &msg.PlayerId)
	return nil
}

type PlayerSelectionErrorResp struct{}

func (msg *PlayerSelectionErrorResp) Opcode() string                 { return "AAE" }
func (msg *PlayerSelectionErrorResp) Serialize(out io.Writer) error  { return nil }
func (msg *PlayerSelectionErrorResp) Deserialize(in io.Reader) error { return nil }

type PlayerSelectionResp struct {
	Player *db.Player
	// TODO items
}

func (msg *PlayerSelectionResp) Opcode() string { return "ASK" }

func (msg *PlayerSelectionResp) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "|%d|%s|%d|%d|%d|%d|%v|%v|%v|",
		msg.Player.Id,
		msg.Player.Name,
		msg.Player.Experience.Level,
		msg.Player.Breed,
		btoi(msg.Player.Gender),
		msg.Player.Appearance.Skin,
		msg.Player.Appearance.Colors.First,
		msg.Player.Appearance.Colors.Second,
		msg.Player.Appearance.Colors.Third,
		// TODO send items
	)
	return nil
}

func (msg *PlayerSelectionResp) Deserialize(in io.Reader) error { return nil }

type ContextType int

const (
	InvalidContextType ContextType = iota
	SoloContextType
	FightContextType
)

type GameContextCreateReq struct {
	Type ContextType
}

func (msg *GameContextCreateReq) Opcode() string { return "GC" }

func (msg *GameContextCreateReq) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%d", int(msg.Type))
	return nil
}

func (msg *GameContextCreateReq) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "GC%d", &msg.Type)
	return nil
}

type GameContextCreateResp struct {
	Type ContextType
}

func (msg *GameContextCreateResp) Opcode() string { return "GCK" }

func (msg *GameContextCreateResp) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "|%d|", int(msg.Type))
	return nil
}

func (msg *GameContextCreateResp) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "GCK|%d|", &msg.Type)
	return nil
}

type GameContextCreateErrorResp struct{}

func (msg *GameContextCreateErrorResp) Opcode() string                 { return "GCE" }
func (msg *GameContextCreateErrorResp) Serialize(out io.Writer) error  { return nil }
func (msg *GameContextCreateErrorResp) Deserialize(in io.Reader) error { return nil }

type SetPlayerStats struct {
	Experience            uint64
	LowerExperienceLevel  uint64
	HigherExperienceLevel uint64

	Kamas uint64

	BoostStatsPts  int
	BoostSpellsPts int

	AlignId    int
	AlignLevel int
	AlignGrade int
	Honor      int
	Dishonor   int
	EnabledPvp bool

	Life    int
	MaxLife int

	Energy    int
	MaxEnergy int

	Stats sdb.StatList
}

func (msg *SetPlayerStats) Opcode() string { return "As" }

func (msg *SetPlayerStats) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%d,%d,%d", msg.Experience, msg.LowerExperienceLevel, msg.HigherExperienceLevel)

	fmt.Fprintf(out, "|%d", msg.Kamas)

	fmt.Fprintf(out, "|%d|%d", msg.BoostStatsPts, msg.BoostSpellsPts)

	fmt.Fprintf(out, "|%d~%d,%d,%d,%d,%d", msg.AlignId, msg.AlignLevel, msg.AlignGrade, msg.Honor, msg.Dishonor, btoi(msg.EnabledPvp))

	fmt.Fprintf(out, "|%d,%d", msg.Life, msg.MaxLife)

	fmt.Fprintf(out, "|%d,%d", msg.Energy, msg.MaxEnergy)

	for _, o := range msg.Stats.Stats() {
		fmt.Fprint(out, "|")

		switch stat := o.(type) {
		case sdb.MinStat:
			fmt.Fprintf(out, "%d", stat.Total())
		case sdb.Stat:
			fmt.Fprintf(out, "%d,%d,%d,%d", stat.Base(), stat.Equipment(), stat.Gift(), stat.Context())
		case sdb.ExStat:
			fmt.Fprintf(out, "%d,%d,%d,%d,%d", stat.Base(), stat.Equipment(), stat.Gift(), stat.Context(), stat.Total())
		}
	}

	return nil
}

func (msg *SetPlayerStats) Deserialize(in io.Reader) error {
	return nil
}
