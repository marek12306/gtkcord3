package completer

import (
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/gtkcord3/gtkcord/md"
	"github.com/diamondburned/gtkcord3/gtkcord/semaphore"
	"github.com/diamondburned/gtkcord3/internal/log"
	"github.com/gotk3/gotk3/gtk"
)

func (c *State) completeEmojis(word string) {
	if word == "" {
		return
	}

	guildID := c.container.GetGuildID()
	guildEmojis, err := c.state.Emoji.Get(guildID)
	if err != nil {
		log.Errorln("Failed to get emojis:", err)
		return
	}

	filtered := guildEmojis[:0]
	filteredLen := 0

	for _, guild := range guildEmojis {
		filteredEmojis := make([]discord.Emoji, 0, MaxCompletionEntries)

		for _, e := range guild.Emojis {
			if contains(e.Name, word) {
				filteredEmojis = append(filteredEmojis, e)
				filteredLen++

				if filteredLen > MaxCompletionEntries {
					break
				}
			}
		}

		guild.Emojis = filteredEmojis
		filtered = append(filtered, guild)

		if filteredLen > MaxCompletionEntries {
			break
		}
	}

	if len(filtered) == 0 {
		return
	}

	semaphore.IdleMust(func() {
		for _, guild := range filtered {
			for _, e := range guild.Emojis {
				b, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

				b.Add(completerImage(md.EmojiURL(e.ID.String(), e.Animated)))
				b.Add(completerLeftLabel(e.Name))
				b.Add(completerRightLabel(guild.Name))
				c.addCompletionEntry(b, e.String())
			}
		}
	})
}
