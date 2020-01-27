package dgperms

import "github.com/bwmarrin/discordgo"

type Permissions = int

const (
	// User may create instant invites
	PermissionCreateInstantInvite = 1 << iota
	// User may kick members
	PermissionKickMembers    // may require 2fa
	PermissionBanMembers     // may require 2fa
	PermissionAdministrator  // may require 2fa
	PermissionManageChannels // may require 2fa
	PermissionManageGuild    // may require 2fa
	PermissionAddReactions
	PermissionViewAuditLog
	PermissionViewChannel
	PermissionSendMessages
	PermissionSendTTSMessages
	PermissionManageMessages // may require 2fa
	PermissionEmbedLinks
	PermissionAttachFiles
	PermissionReadMessageHistory
	PermissionMentionEveryone
	PermissionUseExternalEmojis
	PermissionConnect
	PermissionSpeak
	PermissionMuteMembers
	PermissionDeafenMembers
	PermissionMoveMembers
	PermissionUseVAD
	PermissionPrioritySpeaker
	PermissionStream
	PermissionChangeNickname
	PermissionManageNicknames
	PermissionManageRoles    // may require 2fa
	PermissionManageWebhooks // may require 2fa
	PermissionManageEmojis   // may require 2fa

	PermissionAllText = 0 |
		PermissionViewChannel |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone

	PermissionAllVoice = 0 |
		PermissionConnect |
		PermissionSpeak |
		PermissionMuteMembers |
		PermissionDeafenMembers |
		PermissionMoveMembers |
		PermissionUseVAD |
		PermissionPrioritySpeaker

	PermissionAllChannel = 0 |
		PermissionAllText |
		PermissionAllVoice |
		PermissionCreateInstantInvite |
		PermissionManageRoles |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionViewAuditLog |
		PermissionStream

	PermissionAll = 0 |
		PermissionAllChannel |
		PermissionKickMembers |
		PermissionBanMembers |
		PermissionManageGuild |
		PermissionAdministrator |
		PermissionManageWebhooks |
		PermissionManageEmojis |
		PermissionChangeNickname |
		PermissionManageNicknames |
		PermissionUseExternalEmojis
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
