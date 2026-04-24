import {
  defineConfig,
  presetAttributify,
  presetIcons,
  presetUno,
  transformerDirectives,
  transformerVariantGroup,
} from 'unocss';

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
        ink: 'var(--bb-color-ink)',
        muted: 'var(--bb-color-muted)',
        line: 'var(--bb-color-line)',
        paper: 'var(--bb-color-paper)',
        surface: 'var(--bb-color-surface)',
        honey: 'var(--bb-color-honey)',
        leaf: 'var(--bb-color-leaf)',
        blue: 'var(--bb-color-blue)',
        violet: 'var(--bb-color-violet)',
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
    'bb-scrollbar': 'scrollbar-thin scrollbar-thumb-#c7c3ba scrollbar-track-transparent',
  },
});
