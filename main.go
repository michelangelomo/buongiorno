package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	images []*image
	token string
	tokenEnvVar = "DISCORD_BOT_TOKEN"
)

type image struct {
	baseAddress string
	command string
	images []string
	pages int
} 

func main() {
	images = append(images, &image{baseAddress: "https://ticondivido.it/buongiorno/", command: "buongiorno", pages: 6})
	images = append(images, &image{baseAddress: "https://ticondivido.it/buonanotte-2/", command: "buonanotte", pages: 6})
	
	images = append(images, &image{baseAddress: "https://ticondivido.it/immagini-buon-lunedi/", command: "lunedi", pages: 6})
	images = append(images, &image{baseAddress: "https://ticondivido.it/immagini-buon-martedi/", command: "martedi", pages: 6})
	images = append(images, &image{baseAddress: "https://ticondivido.it/immagini-buon-mercoledi/", command: "mercoledi", pages: 6})
	images = append(images, &image{baseAddress: "https://ticondivido.it/immagini-buon-giovedi/", command: "giovedi", pages: 6})
	images = append(images, &image{baseAddress: "https://ticondivido.it/immagini-buon-venerdi/", command: "venerdi", pages: 6})
	images = append(images, &image{baseAddress: "https://ticondivido.it/immagini-buon-sabato/", command: "sabato", pages: 6})
	images = append(images, &image{baseAddress: "https://ticondivido.it/buona-domenica/", command: "domenica", pages: 6})
	images = append(images, &image{baseAddress: "https://ticondivido.it/buon-weekend/", command: "weekend", pages: 6})
	images = append(images, &image{baseAddress: "https://ticondivido.it/buon-compleanno/", command: "compleanno", pages: 4})
	images = append(images, &image{baseAddress: "https://ticondivido.it/buongiorno-divertenti/", command: "buongiorno-divertente", pages: 3})
	images = append(images, &image{baseAddress: "https://ticondivido.it/immagini-buongiorno-natalizio/", command: "buongiorno-natalizio", pages: 3})
	images = append(images, &image{baseAddress: "https://ticondivido.it/buona-pensione/", command: "pensione", pages: 0})
	images = append(images, &image{baseAddress: "https://ticondivido.it/immagini-festa-della-donna/", command: "donna", pages: 0})
	images = append(images, &image{baseAddress: "https://ticondivido.it/immagini-buon-25-aprile/", command: "25-aprile", pages: 0})
	images = append(images, &image{baseAddress: "https://ticondivido.it/immagini-buon-1-maggio/", command: "1-maggio", pages: 0})
	images = append(images, &image{baseAddress: "https://ticondivido.it/immagini-buon-2-giugno/", command: "2-giugno", pages: 0})
	images = append(images, &image{baseAddress: "https://ticondivido.it/immagini-buon-ferragosto/", command: "ferragosto", pages: 0})

	for _, i := range(images) {
		i.loadImages()
	}

	token = os.Getenv(tokenEnvVar)
	
	initDiscordBot()
}

func (im *image) loadImages() {
	html := retrieve(im.baseAddress, im.pages)
	var myRegex = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	var imgTags = myRegex.FindAllStringSubmatch(string(html), -1)
	out := make([]string, len(imgTags))
	for i := range out {
		if strings.Contains(imgTags[i][1], "logo") {
			continue
		}
		im.images = append(im.images, imgTags[i][1])
	}
}

func retrieve(mainUrl string, page int) string {
	var concatenatedHTML string
	url := mainUrl
	
	for i := 0; i <= page; i++ {
		if i != 0 {
			url = fmt.Sprintf("%s%d", mainUrl, i)
		}
		fmt.Println("Retrieving", url, "...")

		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		html, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		concatenatedHTML += string(html)
	}
	return concatenatedHTML
}

func (im *image) getRandomImage() string {
    return im.images[rand.Intn(len(im.images))]
}

func initDiscordBot() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	for _, image := range(images) {
		if image.command == strings.ReplaceAll(m.Content, "!", "") {
			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Image: &discordgo.MessageEmbedImage{URL: image.getRandomImage()}})
		}
	}
}