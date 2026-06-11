package utils

import (
	"fmt"
	"slices"

	"github.com/bwmarrin/discordgo"
)

// Yetki kontrol sonucu taslağı
type PermissionResult struct {
	OK     bool  // Yetkili mi durumu
	Reason error // Kabul edilmediyse sebebi
}

// (internal) | Başarılı PermissionResult döndürür
func permOK() PermissionResult {
	return PermissionResult{
		OK: true,
	}
}

// (internal) | Başarısız PermissionResult döndürür
func permFail(reason string) PermissionResult {
	return PermissionResult{
		OK:     false,
		Reason: fmt.Errorf("%s", reason),
	}
}

// Kullanıcının belirtilen kanalda istenen yetkiye sahipliğini kontrol eder
// Kullanım:
//
//	result := utils.HasPermissionAtChannel(s, channelID, userID, discordgo.PermissionSendMessages)
func HasPermissionAtChannel(s *discordgo.Session, channelId, userId string, perm int64) PermissionResult {

	fetchedUserPerms, result, ok := getChannelPerms(s, channelId, userId)
	if !ok {
		return result
	}

	if fetchedUserPerms&perm == 0 {
		// Kullanıcının istenen yetkisi yok.
		return permFail("Yetersiz izin")
	}

	return permOK()

}

// Kullanıcının belirtilen kanalda istenen yetkilere sahipliğini kontrol eder
// Kullanım:
//
//	result := utils.HasAnyPermissionAtChannel(s, channelID, userID, discordgo.PermissionSendMessages , ...)
func HasAnyPermissionAtChannel(s *discordgo.Session, channelId, userId string, perms ...int64) PermissionResult {

	fetchedUserPerms, result, ok := getChannelPerms(s, channelId, userId)
	if !ok {
		return result
	}

	for _, p := range perms {
		if fetchedUserPerms&p != 0 {
			// Kullanıcınıda istenen izinlerden en az bir tanesi var
			return permOK()
		}
	}

	return permFail("Yetersiz izin")
}

// Kullanıcının belirtilen kanalda istenen tüm yetkilere sahipliğini kontrol eder
// Kullanım:
//
//	result := utils.HasAllPermissionAtChannel(s, channelID, userID, discordgo.PermissionSendMessages , ...)
func HasAllPermissionAtChannel(s *discordgo.Session, channelId, userId string, perms ...int64) PermissionResult {

	fetchedUserPerms, result, ok := getChannelPerms(s, channelId, userId)
	if !ok {
		return result
	}

	for _, p := range perms {
		if fetchedUserPerms&p == 0 {
			// Kullanıcınıda istenen izinlerden en az bir tanesi yok
			return permFail("Yetersiz izin")
		}
	}

	return permOK()
}

// Kullanıcının belirtilen Guild'de istenen yetkiye sahipliğini kontrol eder
// Kullanım:
//
//	result := utils.HasPermissionAtGuild(s, guildId, userId, discordgo.PermissionSendMessages)
func HasPermissionAtGuild(s *discordgo.Session, guildId, userId string, perm int64) PermissionResult {
	fetchedPerms, result, ok := getGuildPerms(s, guildId, userId)
	if !ok {
		return result
	}

	if fetchedPerms&perm == 0 {
		return permFail("Yetersiz izin")
	}

	return permOK()
}

// Kullanıcının belirtilen Guild'de istenen yetkilere sahipliğini kontrol eder
// Kullanım:
//
//	result := utils.HasAnyPermissionAtGuild(s, guildId, userId, discordgo.PermissionSendMessages, ...)
func HasAnyPermissionAtGuild(s *discordgo.Session, guildId, userId string, perms ...int64) PermissionResult {

	fetchedPerms, result, ok := getGuildPerms(s, guildId, userId)
	if !ok {
		return result
	}

	for _, p := range perms {
		if fetchedPerms&p != 0 {

			return permOK()
		}
	}

	return permFail("Yetersiz izin")
}

// Kullanıcının belirtilen Guild'de istenen tüm yetkilere sahipliğini kontrol eder
// Kullanım:
//
//	result := utils.HasAllPermissionAtGuild(s, guildId, userId, discordgo.PermissionSendMessages, ...)
func HasAllPermissionAtGuild(s *discordgo.Session, guildId, userId string, perms ...int64) PermissionResult {

	fetchedPerms, result, ok := getGuildPerms(s, guildId, userId)
	if !ok {
		return result
	}

	for _, p := range perms {
		if fetchedPerms&p == 0 {
			return permFail("Yetersiz izin")
		}
	}

	return permOK()
}

// Kullanıcının belirtilen Guild'de sahipliğini kontrol eder
// Kullanım:
//
//	err , ok := utils.IsGuildOwner(s, guildID, userID)
func IsGuildOwner(s *discordgo.Session, guildId, userId string) (err error, ok bool) {
	guild, err := s.State.Guild(guildId)
	if err != nil {
		return err, false
	}

	return nil, guild.OwnerID == userId
}

// Kullanıcının belirtilen Guild'de Yönetici yetkisine sahipliğini kontrol eder
// Kullanım:
//
//	err , ok := utils.IsGuildAdmin(s, guildID, userID)
func IsGuildAdmin(s *discordgo.Session, guildId, userId string) (err error, ok bool) {
	guild, err := s.State.Guild(guildId)
	if err != nil {
		return err, false
	}

	// Sunucu sahibi ise zaten yetkilidir
	if guild.OwnerID == userId {
		return nil, true
	}

	result := HasPermissionAtGuild(s, guildId, userId, discordgo.PermissionAdministrator)

	return result.Reason, result.OK
}

// Kullanıcının belirtilen rol ID'sine sahipliğini kontrol eder
// Kullanım:
//
//	result := utils.HasRole(s, guildId, userId, roleId)
func HasRole(s *discordgo.Session, guildId, userId, roleID string) PermissionResult {
	member, err := getMember(s, guildId, userId)
	if err != nil {
		return permFail("Kullanıcı izinleri alınamadı: " + err.Error())
	}

	if !slices.Contains(member.Roles, roleID) {
		return permFail("Yetersiz rol")
	}

	return permOK()
}

// Kullanıcının verilen rol ID'lerinden en az birine sahipliğini kontrol eder
// Kullanım:
//
//	result := utils.HasAnyRole(s, guildId, userId, roleId1, ...)
func HasAnyRole(s *discordgo.Session, guildId, userId string, roleIDs ...string) PermissionResult {
	member, err := getMember(s, guildId, userId)
	if err != nil {
		return permFail("Kullanıcı izinleri alınamadı: " + err.Error())
	}

	for _, roleID := range roleIDs {
		if slices.Contains(member.Roles, roleID) {
			return permOK()
		}
	}

	return permFail("Yetersiz rol")
}

// Kullanıcının verilen tüm rol ID'lerine sahip olup olmadığını kontrol eder.
// Kullanım:
//
//	result := utils.HasAllRoles(s, guildId, userId, roleId1, ...)
func HasAllRoles(s *discordgo.Session, guildId, userId string, roleIDs ...string) PermissionResult {
	member, err := getMember(s, guildId, userId)
	if err != nil {
		return permFail("Kullanıcı izinleri alınamadı: " + err.Error())
	}

	for _, roleID := range roleIDs {
		if !slices.Contains(member.Roles, roleID) {
			return permFail("Yetersiz rol")
		}
	}

	return permOK()
}

// Botun belirtilen Guild'de verilen yetkiye sahipliğini kontrol eder
func BotHasPermissionAtGuild(s *discordgo.Session, guildId string, perm int64) PermissionResult {
	if s.State.User == nil {
		return permFail("Bot henüz hazır değil")
	}

	return HasPermissionAtGuild(s, guildId, s.State.User.ID, perm)
}

// Botun belirtilen kanalda verilen yetkiye sahipliğini kontrol eder
func BotHasPermissionAtChannel(s *discordgo.Session, channelId string, perm int64) PermissionResult {
	if s.State.User == nil {
		return permFail("Bot henüz hazır değil")
	}
	return HasPermissionAtChannel(s, channelId, s.State.User.ID, perm)
}

// Verilen kanal ID'sinin bir metin kanalına (GUILD_TEXT) aitliğini kontrol eder
// Kullanım:
//
//	err , ok := utils.IsTextChannel(s, channelId)
func IsTextChannel(s *discordgo.Session, channelID string) (err error, ok bool) {
	ch, err := s.Channel(channelID)
	if err != nil {
		return err, false
	}

	if ch.Type != discordgo.ChannelTypeGuildText {
		return nil, false
	}

	return nil, true
}

// Verilen kanal ID'sinin bir ses kanalına (GUILD_VOICE) aitliğini kontrol eder
// Kullanım:
//
//	err , ok := utils.IsVoiceChannel(s, channelId)
func IsVoiceChannel(s *discordgo.Session, channelID string) (err error, ok bool) {
	ch, err := s.Channel(channelID)
	if err != nil {
		return err, false
	}

	if ch.Type != discordgo.ChannelTypeGuildVoice {
		return nil, false
	}

	return nil, true
}

// Kanalın NSFW (Not Safe For Work) olarak işaretliliğini kontrol eder
// Kullanım:
//
//	err , ok := utils.IsNSFWChannel(s, channelId)
func IsNSFWChannel(s *discordgo.Session, channelID string) (err error, ok bool) {
	ch, err := s.Channel(channelID)
	if err != nil {
		return err, false
	}

	if !ch.NSFW {
		return nil, false
	}

	return nil, true
}

// Verilen kanal ID'sinin bir thread (konu) olup olmadığını kontrol eder
// Kullanım:
//
//	err , ok := utils.IsNSFWChannel(s, channelId)
func IsThread(s *discordgo.Session, channelID string) (err error, ok bool) {
	ch, err := s.Channel(channelID)
	if err != nil {
		return err, false
	}

	isThread := ch.Type == discordgo.ChannelTypeGuildPublicThread ||
		ch.Type == discordgo.ChannelTypeGuildPrivateThread ||
		ch.Type == discordgo.ChannelTypeGuildNewsThread

	if !isThread {
		return nil, false
	}

	return nil, true
}

// (internal) | Kullanıcıyı önce discord cache'den almayı dener yok ise direkt sorar
func getMember(s *discordgo.Session, guildId, userId string) (*discordgo.Member, error) {
	// Önce cache'e bak
	member, err := s.State.Member(guildId, userId)
	if err == nil {
		return member, nil
	}
	// Cache'de yoksa API'ye git
	return s.GuildMember(guildId, userId)
}

func getChannelPerms(s *discordgo.Session, channelId, userId string) (memberPerm int64, result PermissionResult, ok bool) {
	perms, err := s.State.UserChannelPermissions(userId, channelId)
	if err != nil {
		return 0, permFail("Kullanıcı izinleri alınamadı: " + err.Error()), false
	}
	if perms&discordgo.PermissionAdministrator != 0 {
		return 0, permOK(), false
	}
	return perms, PermissionResult{}, true
}

func getGuildPerms(s *discordgo.Session, guildId, userId string) (memberPerm int64, result PermissionResult, ok bool) {
	member, err := getMember(s, guildId, userId)
	if err != nil {
		return 0, permFail("Kullanıcı izinleri alınamadı: " + err.Error()), false
	}
	if member.Permissions&discordgo.PermissionAdministrator != 0 {
		return 0, permOK(), false
	}
	return member.Permissions, PermissionResult{}, true
}
