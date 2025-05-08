package gemini

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestGemini(t *testing.T) {
	t.Run("Can parse test youtube description 2019 Weekly Track Roundup: 10/6", func(t *testing.T) {
		t.Skip("Skip calling real Gemini API")

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		description := "CHARITY COMPILATION PRE-ORDER:\nhttps://theneedledrop.merchtable.com/music/the-needle-drop-various-artists-vinyl-12\nILRC: https://www.ilrc.org/\n\nPatreon: https://www.patreon.com/theneedledrop\n\nFAV TRACKS Spotify playlist: https://open.spotify.com/user/tndausten/playlist/2zderg88f9HbH54RJBTp1m?si=W8oXCAHvRnSJun4x6VHhdQ\n\nTurntable Lab link: http://turntablelab.com/theneedledrop\n\nAmazon link: http://amzn.to/1KZmdWI\n\n\nSHOUTOUT: CARLY RAE JEPSEN COVERS NO DOUBT\nhttps://open.spotify.com/album/0tUnCfqBLeZlivAkbxvbib\n\n\n!!!BEST TRACKS THIS WEEK!!!\n\nDoja Cat - Bottom Bitch\nhttps://youtu.be/ik0qg-O_2DM\n\nPoppy - I Disagree\nhttps://www.youtube.com/watch?v=6gmswmbosYo&feature=youtu.be\n\nMoor Mother - After Images\nhttps://youtu.be/VeZIqemkrD8\n\nLightning Bolt - Hüsker Dön't\nhttps://lightningbolt.bandcamp.com/track/h-sker-d-nt\n\nNegative Gemini - Bad Baby (Club Mix)\nhttps://youtu.be/_ddbrUq40Iw\n\nAnamanaguchi - Air on Line\nhttps://youtu.be/nnq1ApucY4g\n\nBig Thief - Forgotten Eyes\nhttps://youtu.be/hGD-8f8Wn5M\n\nG.T. - How Dare You\nhttps://youtu.be/rbrdRcwZE6Q\n\nTyler, the Creator - Earfquake (Channel Tres Remix)\nhttps://www.youtube.com/watch?v=T8jx0d9GAF4\n\nDaniel Caesar - CYANIDE REMIX ft. Koffee\nhttps://www.youtube.com/watch?v=mBKXHk2nJ1I\n\nKim Petras - There Will Be Blood\nhttps://www.youtube.com/watch?v=8nBQ8xv2oLY\n\nSunn O))) - Frost (C)\nhttps://youtu.be/Y20qC3qgpps\n\nJacques Greene - For Love\nhttps://youtu.be/GzdMcHhM7tQ\n\nGreat Grandpa - Bloom\nhttps://youtu.be/jFs4Tliyjpg\n\nFloating Points - Anasickmodular\nhttps://youtu.be/Md9gjJlqAxQ\n\nclipping. - Blood on the Fang\nhttps://youtu.be/s9EsHbqmjN4\n\n\n...meh...\n\nRemo Drive - Romeo\nhttps://youtu.be/1DiNlZMBPY0\n\nChromatics - You're No Good\nhttps://youtu.be/PjUblmk4Cyo\n\nGuapdad 4000 - Gucci Pajamas ft. Chance the Rapper & Charlie Wilson\nhttps://www.youtube.com/watch?v=QLw2eTCKaCg\n\nKing Princess - Hit the Back\nhttps://www.youtube.com/watch?v=GyFsbYSajhs\n\nGucci Mane - Big Booty ft. Megan Thee Stallion\nhttps://www.youtube.com/watch?v=b_Kx8tx88oQ\n\nBen Frost - Catastrophic Deliquescence\nhttps://youtu.be/7HNqV3K7di8\n\nBuju Banton - Lend a Hand\nhttps://youtu.be/xUtrvvHre34\n\nSleater-Kinney - ANIMAL\nhttps://youtu.be/pGOO7EE4Lhw\n\nSummer Walker - Playing Games ft. Bryson Tiller\nhttps://www.youtube.com/watch?v=o_6HGBsMHeA\n\nJuice WRLD - Bandit ft. NBA Youngboy\nhttps://www.youtube.com/watch?v=Sw5fNI400E4\n\nTravis Scott - Highest in the Room\nhttps://www.youtube.com/watch?v=tfSS1e3kYeo\nReview: https://www.youtube.com/watch?v=mjVdNIw9LMk\n\nDanny Brown - 3 Tearz ft. Run the Jewels (prod. JPEGMAFIA)\nhttps://www.youtube.com/watch?v=ApJ1_ZliXLQ\nReview: https://www.youtube.com/watch?v=WB625tK_FK0\n\nCHVRCHES - Death Stranding\nhttps://youtu.be/mFGq92BYmt4\n\nEOB - Santa Teresa\nhttps://youtu.be/TG-Od2-OTdg\n\n\n!!!WORST TRACKS THIS WEEK!!!\n\nNONE! Yeah, none!\n\n===================================\nSubscribe: http://bit.ly/1pBqGCN\n\nOfficial site: http://theneedledrop.com\n\nTND Twitter: http://twitter.com/theneedledrop\n\nTND Facebook: http://facebook.com/theneedledrop\n\nSupport TND: http://theneedledrop.com/support\n===================================\n\nY'all know this is just my opinion, right?"

		apiKey := os.Getenv("GEMINI_API_KEY")

		client := NewClient(apiKey)

		client.ParseYoutubeDescription(description)
	})

	t.Run("Can find valid spotify id for clipping even with spelling error", func(t *testing.T) {
		t.Skip("Skip calling real Gemini API")

		// REAL
		// https://open.spotify.com/track/0jv5VgdENAPV7lHtBlsaXE

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		tracks := []ParsedTrack{
			{
				Title:     "Blood on the Fang",
				Artist:    "clipping.",
				SpotifyId: "",
			},
		}

		apiKey := os.Getenv("GEMINI_API_KEY")

		client := NewClient(apiKey)

		client.FindSpotifyUrls(tracks)
	})
}
