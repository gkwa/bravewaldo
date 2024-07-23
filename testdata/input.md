---
filetype: product
test:
- this and that
- test2
x:
  apple:
      pear: null
  "y":
  - a
  - b
---

https://google.com

https://google.com


# Hello, World!

This is a **test** file for our Goldmark AST roundtrip.






## Features

1. Lists
2. *Italic*
3. **Bold**

> Blockquotes are supported too.

| Column 1 | Column 2 |
|----------|----------|
| Cell 1   | Cell 2   |

- [ ] Task 1
- [x] Task 2

Here's some `inline code` and a code block:

```go
func main() {
   fmt.Println("Hello, World!")
}
```

![Seaweed Salad photo](https://static.spotapps.co/spots/a4/3ebb855c2348c68c7b94a4956d9662/full)

---

[OpenAI](https://www.openai.com)

~~strikethrough~~

1. First item
   - Subitem
   - Another subitem
2. Second item

[Link to Headers](#headers)

Term
: Definition

Here's a sentence with a footnote.[^1]

[^1]: This is the footnote.

:smile: :heart: :thumbsup:

When $a \ne 0$, there are two solutions to $(ax^2 + bx + c = 0)$ and they are $$ x = {-b \pm \sqrt{b^2-4ac} \over 2a} $$


```
> [!NOTE]
> This is a note.

> [!WARNING]
> This is a warning.
```


Example:
```
<!-- This content will not appear in the rendered Markdown -->
```
