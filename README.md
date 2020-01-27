# dgperms

Permission checking utilities for [discordgo](https://github.com/bwmarrin/discordgo).

Follows the [official specification](https://discordapp.com/developers/docs/topics/permissions#permissions).

## Usage
```go
package main
import (
    "github.com/bwmarrin/discordgo"
    "github.com/cfi2017/dgperms"
)

func main() {
    var guild *discordgo.Guild
    var member *discordgo.Member
    var channel *discordgo.Channel

    // ...

    perms := dgperms.ComputePermissions(guild, member, channel)
    if dgperms.HasPermission(perms, dgperms.PermissionViewChannel) {
        // do something
    }

}
```