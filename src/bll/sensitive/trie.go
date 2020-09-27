package sensitive

type Node struct {
	isRootNode bool
	isPathEnd  bool
	Character  rune
	Children   map[rune]*Node
}

func (this *Node) IsPathEnd() bool {
	return this.isPathEnd
}

func NewNode(character rune) *Node {
	return &Node{
		Character: character,
		Children:  make(map[rune]*Node, 0),
	}
}

func NewRootNode(character rune) *Node {
	return &Node{
		isRootNode: true,
		Character:  character,
		Children:   make(map[rune]*Node, 0),
	}
}

// ------------------------------------------------------------------------------------
type Trie struct {
	Root *Node
}

func NewTrie() *Trie {
	return &Trie{
		Root: NewRootNode(0),
	}
}

func (this *Trie) Add(words ...string) {
	for _, word := range words {
		this.add(word)
	}
}

func (this *Trie) add(word string) {
	var current = this.Root
	var runes = []rune(word)
	for position := 0; position < len(runes); position++ {
		r := runes[position]
		if next, ok := current.Children[r]; ok {
			current = next
		} else {
			newNode := NewNode(r)
			current.Children[r] = newNode
			current = newNode
		}
		if position == len(runes)-1 {
			current.isPathEnd = true
		}
	}
}

func (this *Trie) Replace(text string, character rune) string {
	var (
		parent  = this.Root
		current *Node
		runes   = []rune(text)
		length  = len(runes)
		left    = 0
		found   bool
	)

	for position := 0; position < len(runes); position++ {
		current, found = parent.Children[runes[position]]

		if !found || (!current.IsPathEnd() && position == length-1) {
			parent = this.Root
			position = left
			left++
			continue
		}

		if current.IsPathEnd() && left <= position {
			for i := left; i <= position; i++ {
				runes[i] = character
			}
		}

		parent = current
	}

	return string(runes)
}
