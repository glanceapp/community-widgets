# Figure Card

A simple image widget with an optional caption, internal frame, and padding controls.

## Widget Type

`custom-api`

## Preview

![Figure Card displaying a padded Moe Counter with a caption](./preview.png)

*The preview shows one example using a padded Moe Counter. Frame, padding, and caption behavior can be configured independently.*

## Options

| Name        | Type    | Required | Default |
| ----------- | ------- | -------- | ------- |
| `src`       | string  | Yes      | —       |
| `alt`       | string  | No       | `""`    |
| `caption`   | string  | No       | —       |
| `frameless` | boolean | No       | `false` |
| `padding`   | boolean | No       | `false` |

### `src`

Direct image URL or browser-accessible path, such as `/assets/image.png`.

The browser loads the image directly; Glance does not fetch it through the `custom-api` widget. The widget-level `cache` property therefore does not control image caching, which depends on the image server's HTTP cache headers.

### `alt`

Alternative text for the image. This maps directly to the HTML `alt` attribute and is used for accessibility.

This option may be omitted. The widget then renders `alt=""`.

### `caption`

Optional text displayed below the image.

### `frameless`

Controls the *Figure Card's* internal frame.

This is separate from Glance's widget-level `frameless` property.

* `false`: render an inner Glance-style card frame.
* `true`: remove the inner frame.

### `padding`

Adds internal padding around the image and caption.

Use this for transparent PNGs, counters, diagrams, or other images that should not touch the card edge.

## Configuration

When adding multiple *Figure Cards*, define the template once with a YAML anchor and reuse it with an alias. This keeps the configuration compact and ensures every *Figure Card* uses the same rendering logic.

```yaml
- type: custom-api
  frameless: true
  title: Figure
  options:
    src: https://placehold.co/600x400
    alt: Example alt
    caption: Example caption
    frameless: false
    padding: true
  template: &figure-card-template |
    {{ $src := .Options.StringOr "src" "" }}
    {{ $alt := .Options.StringOr "alt" "" }}
    {{ $caption := .Options.StringOr "caption" "" }}
    {{ $frameless := .Options.BoolOr "frameless" false }}
    {{ $padding := .Options.BoolOr "padding" false }}

    {{ if $src }}
      <figure
        class="card{{ if not $frameless }} widget-content-frame{{ else }} rounded{{ end }}"
        style="
          margin: 0;
          overflow: hidden;
          {{ if $padding }}padding: var(--widget-content-padding);{{ end }}
        "
      >
        <img
          {{ if or $frameless $padding }}class="rounded"{{ end }}
          loading="lazy"
          decoding="async"
          src="{{ $src }}"
          alt="{{ $alt }}"
          style="display: block; width: 100%; height: auto;"
        >

        {{ with $caption }}
          <figcaption
            class="
              margin-top-10
              flex
              flex-column
              grow
              color-subdue
              {{ if not $padding }}padding-inline-widget margin-bottom-widget{{ end }}
            "
          >
            {{ . }}
          </figcaption>
        {{ end }}
      </figure>
    {{ end }}
```

> **Note:** Set the outer Glance widget to `frameless: true`. Figure Card manages its own internal frame and padding through `options.frameless` and `options.padding`.
