package main
// Importing required packages
import (
	"fmt"
	"strings"
	"unicode"
)

// NodeType defines the type of an HTML node.
type NodeType int

const (
	ElementNode NodeType = iota
	TextNode
)

// Node represents an HTML node, which can be an element or text.
type Node struct {
	Type       NodeType
	Tag        string
	Attributes map[string]string
	Children   []*Node
	Content    string
}

// Parser holds the HTML data and the position of the current character.
type Parser struct {
	html string
	pos  int
}

// NewParser creates a new Parser instance.
func NewParser(html string) *Parser {
	return &Parser{html: html, pos: 0}
}

// Parse parses the HTML and returns the root node.
func (p *Parser) Parse() (*Node, error) {
	return p.parseNodes()
}

// parseNodes parses a list of nodes.
func (p *Parser) parseNodes() (*Node, error) {
	var root Node
	root.Type = ElementNode
	root.Tag = "root"

	for p.pos < len(p.html) {
		p.skipWhitespace()
		if p.pos < len(p.html) && p.html[p.pos] == '<' {
			if p.html[p.pos+1] == '/' {
				p.pos += 2 // skip '</'
				p.parseTagName() // read until '>'
				p.pos++          // skip '>'
				break
			} else {
				child, err := p.parseElement()
				if err != nil {
					return nil, err
				}
				root.Children = append(root.Children, child)
			}
		} else {
			textNode := p.parseText()
			if textNode != nil {
				root.Children = append(root.Children, textNode)
			}
		}
	}
	return &root, nil
}

// parseElement parses an HTML element.
func (p *Parser) parseElement() (*Node, error) {
	p.pos++ // skip '<'
	tag := p.parseTagName()
	attributes := p.parseAttributes()

	node := &Node{
		Type:       ElementNode,
		Tag:        tag,
		Attributes: attributes,
	}

	if p.pos < len(p.html) && p.html[p.pos] == '>' {
		p.pos++ // skip '>'

		// Parse children nodes
		childNodes, err := p.parseNodes()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, childNodes.Children...)
	}
	return node, nil
}

// parseTagName reads the tag name from the HTML string.
func (p *Parser) parseTagName() string {
	start := p.pos
	for p.pos < len(p.html) && (unicode.IsLetter(rune(p.html[p.pos])) || unicode.IsDigit(rune(p.html[p.pos]))) {
		p.pos++
	}
	return p.html[start:p.pos]
}

// parseAttributes reads attributes of an HTML tag.
func (p *Parser) parseAttributes() map[string]string {
	attributes := make(map[string]string)
	for p.pos < len(p.html) && p.html[p.pos] != '>' {
		p.skipWhitespace()
		name := p.parseTagName()
		var value string
		if p.pos < len(p.html) && p.html[p.pos] == '=' {
			p.pos++ // skip '='
			value = p.parseAttributeValue()
		}
		attributes[name] = value
	}
	return attributes
}

// parseAttributeValue reads the attribute value enclosed in quotes.
func (p *Parser) parseAttributeValue() string {
	if p.pos >= len(p.html) || p.html[p.pos] != '"' {
		return ""
	}
	p.pos++ // skip opening quote
	start := p.pos
	for p.pos < len(p.html) && p.html[p.pos] != '"' {
		p.pos++
	}
	value := p.html[start:p.pos]
	p.pos++ // skip closing quote
	return value
}

// parseText parses a text node.
func (p *Parser) parseText() *Node {
	start := p.pos
	for p.pos < len(p.html) && p.html[p.pos] != '<' {
		p.pos++
	}
	content := strings.TrimSpace(p.html[start:p.pos])
	if content == "" {
		return nil
	}
	return &Node{Type: TextNode, Content: content}
}

// skipWhitespace skips whitespace characters.
func (p *Parser) skipWhitespace() {
	for p.pos < len(p.html) && unicode.IsSpace(rune(p.html[p.pos])) {
		p.pos++
	}
}

// printNode recursively prints the HTML structure.
func printNode(node *Node, indent int) {
	indentation := strings.Repeat("  ", indent)
	if node.Type == ElementNode {
		fmt.Printf("%s<%s", indentation, node.Tag)
		for attr, val := range node.Attributes {
			fmt.Printf(" %s=\"%s\"", attr, val)
		}
		fmt.Println(">")
		for _, child := range node.Children {
			printNode(child, indent+1)
		}
		fmt.Printf("%s</%s>\n", indentation, node.Tag)
	} else if node.Type == TextNode {
		fmt.Printf("%s%s\n", indentation, node.Content)
	}
}

func main() {
	html := `<html>
	<head><title>=HTML Parser</title></head>
	<body>
		<h1>Welcome to the Sample Page</h1>
		<p>This is a HTML parser.</p>
	</body>
</html>`

	parser := NewParser(html)
	root, err := parser.Parse()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	printNode(root, 0)
}
