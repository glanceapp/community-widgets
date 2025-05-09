## Governance

**These rules are subject to change as the repository grows and the community evolves.**

* Everyone can contribute widgets to this repository
* Duplicate widget types are allowed so long as they are sufficiently different from one another, so for example you could have two or more different weather widgets that use different APIs, provide different information or have different styles
* Pull requests require at least one vouch from the community before being merged so long as they meet the guidelines
* Anyone can vouch for a pull request by leaving a comment with a 👍 or anything of the like after having tested the widget, though it must be a comment and not a reaction
* The authors of each widget are responsible for maintaining it, including updating it to work with new versions of Glance
* The authors of each widget are responsible for responding to issues and pull requests related to their widget
* If a widget is left in a broken state for an extended period of time and with no response from its author, it may be removed from the repository or transferred to a new maintainer
* If you are submitting a pull request that modifies the widget of another author, please mention them in the pull request description and give them adequate time to review it
* If you no longer wish to be the maintainer of a widget, please submit an issue

## Submitting your `extension` widget

1. Fork the repository
2. Add your widget to the list in the `README.md` file in the root of the repository in alphabetical order using the format `* [Widget name](link to your repository) by @your-username - description`
3. Commit your changes and push them to your fork, then create a pull request

## Submitting your `custom-api` widget

If you have created a custom widget that you would like to share with the community you can do so by following these steps:

1. Fork the repository
2. Create a new directory under `widgets` with the name of your widget and your GitHub username (e.g. `my-widget-by-my-username`)
3. Make sure to read and follow the guidelines below
4. Create a `README.md` file in the new directory and include a preview of the widget, its YAML content and any relevant information required to configure the widget, checkout [this](widgets/immich-stats-by-svilenmarkov/README.md) widget as an example
5. Place your preview image in the same directory with the name `preview.png`
6. Add your widget to the list in the `README.md` file in the root of the repository in alphabetical order using the format `* [Widget name](widgets/widget-directory/README.md) by @your-username - description`
7. Commit your changes and push them to your fork, then create a pull request

### Guidelines

### Add sufix's to your css styles
If you want to use custom css styling make sure to add suffix to your classes based on the name of your widget as shown below:

```css
.{class}-{widget_name} {
  text-align: center;
  margin-top: 1.5rem;
}
```

#### Use descriptive environment variables in place of configurable values
Don't submit local domains or IP's:

```yaml
url: https://192.168.0.50:8080/api/server/statistics
```

Use a variable instead:

```yaml
url: https://${IMMICH_URL}/api/server/statistics
```

<hr>

Don't submit API keys:

```yaml
headers:
  x-api-key: 1234567890
```

Use a variable instead:

```yaml
headers:
  x-api-key: ${IMMICH_API_KEY}
```

#### Provide reasonable default cache times

This can be a few minutes or less for requests to local services or multiple hours for external data that doesn't change frequently.

#### Don't submit widgets which require additional local API's

The goal of this repository is to provide widgets which can be copy-pasted, slightly configured and work out of the box. If your widget requires an additional API to run alongside it that does some kind of data parsing, please create an extension widget instead or request new functionality in the [main repository](https://github.com/glanceapp/glance) for the `custom-api` widget that would allow that parsing to be done directly in the widget.

#### Apply the necessary custom style directly to the elements

Don't rely on the CSS classes that are specific to existing widgets:

```html
<img class="twitch-category-thumbnail" src="...">
```

Instead, embed the necessary CSS directly:

```html
<img style="width: 5rem; aspect-ratio: 3 / 4; border-radius: var(--border-radius);" src="...">
```

This reduces the odds of your widget breaking in the future if the CSS classes get modified. You can use all the utility classes such as `flex`, `color-primary`, `size-h3`, `text-center`, etc, as those are unlikely to change. If you have an idea for a new utility class, please open an issue or submit a PR in the [main repository](https://github.com/glanceapp/glance).

#### Don't hardcode colors

Use the `color-primary`, `color-positive`, `color-negative`, `color-highlight`, `color-subdue`, etc utility classes instead.

#### You can submit multiple different variants of the same widget

If you want to provide multiple styles or layouts for the same widget, you can include them in the same `README.md` file, just make sure to include a preview image for each variant. You can also include multiple variants of the same widget to account for different API versions.
