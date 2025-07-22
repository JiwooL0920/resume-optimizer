// Performance utilities for the Resume Optimizer application

/**
 * Debounce function to limit the rate of function execution
 * @param func - Function to debounce
 * @param wait - Wait time in milliseconds
 * @param immediate - Execute immediately on first call
 */
export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  wait: number,
  immediate?: boolean
): ((...args: Parameters<T>) => void) => {
  let timeout: NodeJS.Timeout | null = null
  
  return function executedFunction(...args: Parameters<T>) {
    const later = () => {
      timeout = null
      if (!immediate) func(...args)
    }
    
    const callNow = immediate && !timeout
    
    if (timeout) clearTimeout(timeout)
    timeout = setTimeout(later, wait)
    
    if (callNow) func(...args)
  }
}

/**
 * Throttle function to limit the rate of function execution
 * @param func - Function to throttle
 * @param limit - Time limit in milliseconds
 */
export const throttle = <T extends (...args: any[]) => any>(
  func: T,
  limit: number
): ((...args: Parameters<T>) => void) => {
  let inThrottle: boolean
  
  return function executedFunction(...args: Parameters<T>) {
    if (!inThrottle) {
      func(...args)
      inThrottle = true
      setTimeout(() => (inThrottle = false), limit)
    }
  }
}

/**
 * Lazy loading intersection observer for images and components
 * @param callback - Function to call when element intersects
 * @param options - IntersectionObserver options
 */
export const createLazyLoader = (
  callback: (entry: IntersectionObserverEntry) => void,
  options: IntersectionObserverInit = {}
): IntersectionObserver => {
  const defaultOptions: IntersectionObserverInit = {
    root: null,
    rootMargin: '50px',
    threshold: 0.1,
    ...options,
  }
  
  return new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        callback(entry)
      }
    })
  }, defaultOptions)
}

/**
 * Preload images for better performance
 * @param imageUrls - Array of image URLs to preload
 */
export const preloadImages = (imageUrls: string[]): Promise<void[]> => {
  const promises = imageUrls.map((url) => {
    return new Promise<void>((resolve, reject) => {
      const img = new Image()
      img.onload = () => resolve()
      img.onerror = () => reject(new Error(`Failed to load image: ${url}`))
      img.src = url
    })
  })
  
  return Promise.all(promises)
}

/**
 * Memoization utility for expensive computations
 * @param fn - Function to memoize
 * @param keyGenerator - Function to generate cache key
 */
export const memoize = <T extends (...args: any[]) => any>(
  fn: T,
  keyGenerator?: (...args: Parameters<T>) => string
): T => {
  const cache = new Map<string, ReturnType<T>>()
  
  return ((...args: Parameters<T>): ReturnType<T> => {
    const key = keyGenerator ? keyGenerator(...args) : JSON.stringify(args)
    
    if (cache.has(key)) {
      return cache.get(key)!
    }
    
    const result = fn(...args)
    cache.set(key, result)
    
    return result
  }) as T
}

/**
 * Virtual scrolling helper for large lists
 */
export class VirtualScroller {
  private container: HTMLElement
  private itemHeight: number
  private bufferSize: number
  private scrollTop: number = 0
  private visibleStart: number = 0
  private visibleEnd: number = 0
  
  constructor(container: HTMLElement, itemHeight: number, bufferSize: number = 5) {
    this.container = container
    this.itemHeight = itemHeight
    this.bufferSize = bufferSize
  }
  
  calculateVisibleRange(totalItems: number): { start: number; end: number } {
    const containerHeight = this.container.clientHeight
    const visibleItemCount = Math.ceil(containerHeight / this.itemHeight)
    
    this.visibleStart = Math.max(0, Math.floor(this.scrollTop / this.itemHeight) - this.bufferSize)
    this.visibleEnd = Math.min(totalItems, this.visibleStart + visibleItemCount + this.bufferSize * 2)
    
    return {
      start: this.visibleStart,
      end: this.visibleEnd
    }
  }
  
  updateScrollTop(scrollTop: number): void {
    this.scrollTop = scrollTop
  }
  
  getOffsetTop(): number {
    return this.visibleStart * this.itemHeight
  }
  
  getTotalHeight(totalItems: number): number {
    return totalItems * this.itemHeight
  }
}

/**
 * Bundle size analyzer for development
 */
export const analyzeBundle = () => {
  if (process.env.NODE_ENV !== 'development') return
  
  console.group('📦 Bundle Analysis')
  
  // Analyze loaded scripts
  const scripts = Array.from(document.querySelectorAll('script[src]'))
  const totalScripts = scripts.length
  
  console.log(`📄 Total Scripts: ${totalScripts}`)
  
  // Analyze CSS
  const stylesheets = Array.from(document.querySelectorAll('link[rel="stylesheet"]'))
  console.log(`🎨 Total Stylesheets: ${stylesheets.length}`)
  
  // Memory usage (if available)
  if ('memory' in performance) {
    const memory = (performance as any).memory
    console.log(`🧠 Memory Usage:`, {
      used: `${Math.round(memory.usedJSHeapSize / 1024 / 1024)} MB`,
      total: `${Math.round(memory.totalJSHeapSize / 1024 / 1024)} MB`,
      limit: `${Math.round(memory.jsHeapSizeLimit / 1024 / 1024)} MB`
    })
  }
  
  console.groupEnd()
}

/**
 * Performance monitor for React components
 */
export const createPerformanceMonitor = (componentName: string) => {
  return {
    startRender: () => {
      if (process.env.NODE_ENV === 'development') {
        performance.mark(`${componentName}-render-start`)
      }
    },
    
    endRender: () => {
      if (process.env.NODE_ENV === 'development') {
        performance.mark(`${componentName}-render-end`)
        performance.measure(
          `${componentName}-render`,
          `${componentName}-render-start`,
          `${componentName}-render-end`
        )
        
        const measures = performance.getEntriesByName(`${componentName}-render`)
        const lastMeasure = measures[measures.length - 1]
        if (lastMeasure && lastMeasure.duration > 16) {
          console.warn(`⚠️ ${componentName} render took ${lastMeasure.duration.toFixed(2)}ms`)
        }
      }
    }
  }
}

/**
 * Resource hints for better loading performance
 */
export const addResourceHints = () => {
  const head = document.head
  
  // Preconnect to external domains
  const domains = [
    'https://fonts.googleapis.com',
    'https://fonts.gstatic.com',
  ]
  
  domains.forEach(domain => {
    const link = document.createElement('link')
    link.rel = 'preconnect'
    link.href = domain
    head.appendChild(link)
  })
}

/**
 * Web Workers utility for offloading heavy computations
 */
export const createWorker = (workerFunction: Function): Worker => {
  const blob = new Blob([`(${workerFunction.toString()})()`], {
    type: 'application/javascript'
  })
  
  return new Worker(URL.createObjectURL(blob))
}