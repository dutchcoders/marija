package mention

import (
	"fmt"
	"strings"
)

func ExampleGetTags_mention() {
	msg := " hello @gernest"
	tags := GetTags('@', strings.NewReader(msg))
	fmt.Println(tags)

	//Output:
	//[gernest]
}

func ExampleGetTags_hashtag() {
	msg := " viva la #tanzania"
	tags := GetTags('#', strings.NewReader(msg))
	fmt.Println(tags)

	//Output:
	//[tanzania]
}
