package htmlutils

import (
	"bytes"
	"errors"
	"golang.org/x/net/html"
	"io"
	"strings"
)

func QuerySelector(doc *html.Node, tag string, attr string, val string) (nodes []*html.Node, err error) {
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tag {
			for _, a := range node.Attr {
				if a.Key == attr {
					if strings.Contains(a.Val, val) == true {
						nodes = append(nodes, node)
						return
					}
				}

			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if nodes != nil {
		return nodes, nil
	}
	return nil, errors.New("Missing " + tag + " with " + attr + " = " + val)
}

func GetGeneralTags(doc *html.Node, tag string) (nodes []*html.Node, err error) {
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tag {
				nodes = append(nodes, node)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if nodes != nil {
		return nodes, nil
	}
	return nil, errors.New("Missing" + tag + "in the node tree")
}

func GetNodeText(node *html.Node, tag string) (nodes []byte) {

	doc := strings.NewReader(RenderNode(node))

	z := html.NewTokenizer(doc)

	depth := 0
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.TextToken:
			if depth > 0 {
				// emitBytes should copy the []byte it receives,
				// if it doesn't process it immediately.
				return z.Text()
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
				if bytes.Equal(tn, []byte(tag)) == true {
					if tt == html.StartTagToken {
						depth++
					} else {
						depth--
					}
				}
				if tt == html.StartTagToken {
					depth++
				} else {
					depth--
				}
		}
	}
}

func GetValueAttr(doc *html.Node, tag string, attr string) (nodes [][]byte, err error) {
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tag {
			for _, a := range node.Attr {
				if a.Key == attr {
					//nodes = append(nodes, []byte(a.Val)...)
					//nodes = append(nodes, []byte("\n"))
					nodes = append(nodes, []byte(a.Val))
					return
				}

			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if nodes != nil {
		return nodes, nil
	}
	return nil, errors.New("Missing \"value\" in the attribute tag")
}

func RenderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}
