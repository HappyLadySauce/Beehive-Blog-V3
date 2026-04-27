import {
  defineConfig,
  presetAttributify,
  presetIcons,
  presetUno,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss'

export default defineConfig({
  presets: [
    presetUno(),
    presetAttributify(),
    presetIcons({
      scale: 1.05,
      warn: true,
    }),
  ],
  transformers: [transformerDirectives(), transformerVariantGroup()],
  theme: {
    breakpoints: {
      sm: '640px',
      md: '768px',
      lg: '1024px',
      xl: '1280px',
      '2xl': '1536px',
    },
    colors: {
      brand: {
        bg: 'var(--bb-color-bg)',
        surface: 'var(--bb-color-surface)',
        subtle: 'var(--bb-color-subtle)',
        text: 'var(--bb-color-text)',
        muted: 'var(--bb-color-muted)',
        line: 'var(--bb-color-line)',
        primary: 'var(--bb-color-primary)',
        primaryStrong: 'var(--bb-color-primary-strong)',
        accent: 'var(--bb-color-accent)',
        danger: 'var(--bb-color-danger)',
        warning: 'var(--bb-color-warning)',
        success: 'var(--bb-color-success)',
      },
    },
    borderRadius: {
      sm: '4px',
      md: '6px',
      lg: '8px',
    },
    boxShadow: {
      panel: 'var(--bb-shadow-panel)',
      focus: '0 0 0 3px var(--bb-color-focus)',
    },
  },
  shortcuts: {
    'bb-focus': 'focus-visible:outline-none focus-visible:shadow-focus',
    'bb-panel': 'border border-brand-line bg-brand-surface shadow-panel rounded-lg',
    'bb-link': 'text-brand-primary underline-offset-4 hover:underline bb-focus rounded-sm',
  },
})
