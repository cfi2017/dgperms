package dgperms

import "github.com/bwmarrin/discordgo"

type Permissions = int

const (
	// User may create instant invites
	PermissionCreateInstantInvite = 0x00000001
	// User may kick members
	PermissionKickMembers        = 0x00000002 // may require 2fa
	PermissionBanMembers         = 0x00000004 // may require 2fa
	PermissionAdministrator      = 0x00000008 // may require 2fa
	PermissionManageChannels     = 0x00000010 // may require 2fa
	PermissionManageGuild        = 0x00000020 // may require 2fa
	PermissionAddReactions       = 0x00000040
	PermissionViewAuditLog       = 0x00000080
	PermissionViewChannel        = 0x00000400
	PermissionSendMessages       = 0x00000800
	PermissionSendTTSMessages    = 0x00001000
	PermissionManageMessages     = 0x00002000 // may require 2fa
	PermissionEmbedLinks         = 0x00004000
	PermissionAttachFiles        = 0x00008000
	PermissionReadMessageHistory = 0x00010000
	PermissionMentionEveryone    = 0x00020000
	PermissionUseExternalEmojis  = 0x00040000
	PermissionConnect            = 0x00100000
	PermissionSpeak              = 0x00200000
	PermissionMuteMembers        = 0x00400000
	PermissionDeafenMembers      = 0x00800000
	PermissionMoveMembers        = 0x01000000
	PermissionUseVAD             = 0x02000000
	PermissionPrioritySpeaker    = 0x00000100
	PermissionStream             = 0x00000200
	PermissionChangeNickname     = 0x04000000
	PermissionManageNicknames    = 0x08000000
	PermissionManageRoles        = 0x10000000 // may require 2fa
	PermissionManageWebhooks     = 0x20000000 // may require 2fa
	PermissionManageEmojis       = 0x40000000 // may require 2fa
	PermissionAll                = 2146958847
)

func ComputeBasePermissions(guild *discordgo.Guild, member *discordgo.Member) (perms Permissions) {
	if member.User.ID == guild.OwnerID {
		return PermissionAll
	}
	for _, role := range guild.Roles {
		if role.ID == guild.ID { // @everyone
			perms |= role.Permissions
			break
		}
	}
	for _, role := range guild.Roles {
		for _, id := range member.Roles {
			if role.ID == id {
				perms |= role.Permissions
				break
			}
		}
	}
	return perms
}

func ComputePermissionOverrides(base Permissions, member *discordgo.Member, channel *discordgo.Channel) (perms Permissions) {
	if base&PermissionAdministrator == PermissionAdministrator {
		return PermissionAll
	}
	perms = base
	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.ID == channel.GuildID {
			perms |= overwrite.Allow
			perms &^= overwrite.Deny
			break
		}
	}

	var allow, deny Permissions
	for _, overwrite := range channel.PermissionOverwrites {
		for _, id := range member.Roles {
			if overwrite.ID == id {
				allow |= overwrite.Allow
				deny |= overwrite.Deny
				break
			}
		}
	}

	perms |= allow
	perms &^= deny

	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.ID == member.User.ID {
			perms |= overwrite.Allow
			perms &^= overwrite.Deny
		}
	}
	return perms
}

func ComputePermissions(guild *discordgo.Guild, member *discordgo.Member, channel *discordgo.Channel) Permissions {
	base := ComputeBasePermissions(guild, member)
	return ComputePermissionOverrides(base, member, channel)
}

func HasPermission(permissions Permissions, permission int) bool {
	return permissions&permission == permission
}
