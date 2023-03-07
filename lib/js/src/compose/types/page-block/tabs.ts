import { PageBlock, PageBlockInput, Registry } from './base'

const kind = 'Tabs'

interface Style {
  appearance: string;
  alignment: string;
  fillJustify: string;
  orientation: string;
  position: string;
}

interface Tab {
  blockID: string;
  title: string;
}

interface Options {
  style: Style;
  tabs: Tab[];
}

const defaults: Readonly<Options> = Object.freeze({
  style: {
    appearance: 'tabs',
    alignment: 'center',
    fillJustify: 'justify',
    orientation: 'horizontal',
    position: 'start',
  },
  tabs: [],
})

export class PageBlockTab extends PageBlock {
  readonly kind = kind

  options: Options = { ...defaults }

  constructor (i?: PageBlockInput) {
    super(i)
    this.applyOptions(i?.options as Partial<Options>)
  }

  applyOptions (o?: Partial<Options>): void {
    if (!o) return

    if (o.tabs) {
      this.options.tabs = o.tabs
    }

    if (o.style) {
      this.options.style = { ...this.options.style, ...o.style }
    }
  }
}

Registry.set(kind, PageBlockTab)
