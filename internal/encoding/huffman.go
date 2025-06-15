package encoding

import (
	"sort"
)

type HuffmanNode struct {
    Char   rune
    Freq   int
    Left   *HuffmanNode
    Right  *HuffmanNode
}

type HuffmanCode struct {
    Codes     map[rune]string
    Tree      *HuffmanNode
    Encoded   string
}

func NewHuffmanCode(content string) *HuffmanCode {
    freqs := make(map[rune]int)
    for _, char := range content {
        freqs[char]++
    }

    var nodes []*HuffmanNode
    for char, freq := range freqs {
        nodes = append(nodes, &HuffmanNode{
            Char: char,
            Freq: freq,
        })
    }

    tree := buildHuffmanTree(nodes)

    codes := make(map[rune]string)
    generateCodes(tree, "", codes)

    encoded := ""
    for _, char := range content {
        encoded += codes[char]
    }

    return &HuffmanCode{
        Codes:   codes,
        Tree:    tree,
        Encoded: encoded,
    }
}

func buildHuffmanTree(nodes []*HuffmanNode) *HuffmanNode {
    for len(nodes) > 1 {
        sort.Slice(nodes, func(i, j int) bool {
            return nodes[i].Freq < nodes[j].Freq
        })

        left := nodes[0]
        right := nodes[1]
        nodes = nodes[2:]


        newNode := &HuffmanNode{
            Freq:  left.Freq + right.Freq,
            Left:  left,
            Right: right,
        }
        nodes = append(nodes, newNode)
    }
    return nodes[0]
}

func generateCodes(node *HuffmanNode, code string, codes map[rune]string) {
    if node == nil {
        return
    }

    if node.Left == nil && node.Right == nil {
        codes[node.Char] = code
        return
    }


    generateCodes(node.Left, code+"0", codes)
    generateCodes(node.Right, code+"1", codes)
}

func (h *HuffmanCode) Decode() string {
    if h.Tree == nil {
        return ""
    }

    result := ""
    current := h.Tree
    
    for _, bit := range h.Encoded {
        if bit == '0' {
            current = current.Left
        } else {
            current = current.Right
        }

        if current.Left == nil && current.Right == nil {
            result += string(current.Char)
            current = h.Tree
        }
    }

    return result
}