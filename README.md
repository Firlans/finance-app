# frontend

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VS Code](https://code.visualstudio.com/) + [Vue (Official)](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Recommended Browser Setup

- Chromium-based browsers (Chrome, Edge, Brave, etc.):
  - [Vue.js devtools](https://chromewebstore.google.com/detail/vuejs-devtools/nhdogjmejiglipccpnnnanhbledajbpd)
  - [Turn on Custom Object Formatter in Chrome DevTools](http://bit.ly/object-formatters)
- Firefox:
  - [Vue.js devtools](https://addons.mozilla.org/en-US/firefox/addon/vue-js-devtools/)
  - [Turn on Custom Object Formatter in Firefox DevTools](https://fxdx.dev/firefox-devtools-custom-object-formatters/)

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Compile and Minify for Production

```sh
npm run build
```

### Run Unit Tests with [Vitest](https://vitest.dev/)

```sh
npm run test:unit
```

### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```

## Component Documentation

Finance app currently exposes two local base components under `src/components/base`.

### BaseRoll

Accordion-like panel wrapper that can render content vertically or horizontally.

#### Props

| Prop | Type | Default | Description |
|---|---|---|---|
| `type` | `String` | `'vertical'` | Accepted values: `horizontal`, `vertical` |
| `label` | `String` | `'Panel'` | Header title |
| `description` | `String` | `''` | Secondary text shown below the label |
| `initiallyOpen` | `Boolean` | `true` | Initial open state |
| `gapClass` | `String` | `'gap-3'` | Gap utility class for slot content |

#### Slots

| Slot | Description |
|---|---|
| default | Main panel content |
| `description` | Replaces the plain description text |

#### Behavior

- `horizontal` mode collapses width and rotates the chevron sideways.
- `vertical` mode collapses height and adds top padding when opened.

### RollBase

Compatibility wrapper around `BaseRoll` with the same public API. Use this component when older templates still reference `RollBase`.

#### Props

`RollBase` forwards the following props directly to `BaseRoll`:

- `type`
- `label`
- `description`
- `initiallyOpen`
- `gapClass`

#### Slots

- default
- `description`
