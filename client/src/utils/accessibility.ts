// Accessibility utilities for the Resume Optimizer application

/**
 * Announces text to screen readers
 * @param message - Text to announce
 * @param priority - Announcement priority ('polite' | 'assertive')
 */
export const announceToScreenReader = (message: string, priority: 'polite' | 'assertive' = 'polite') => {
  const announcement = document.createElement('div')
  announcement.setAttribute('aria-live', priority)
  announcement.setAttribute('aria-atomic', 'true')
  announcement.setAttribute('class', 'sr-only')
  announcement.textContent = message
  
  document.body.appendChild(announcement)
  
  // Remove the announcement after a short delay
  setTimeout(() => {
    document.body.removeChild(announcement)
  }, 1000)
}

/**
 * Focuses the first focusable element within a container
 * @param container - Container element to search within
 */
export const focusFirstElement = (container: HTMLElement) => {
  const focusableElements = container.querySelectorAll(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  )
  
  const firstElement = focusableElements[0] as HTMLElement
  if (firstElement) {
    firstElement.focus()
  }
}

/**
 * Traps focus within a container (useful for modals)
 * @param container - Container to trap focus within
 */
export const trapFocus = (container: HTMLElement) => {
  const focusableElements = container.querySelectorAll(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  ) as NodeListOf<HTMLElement>
  
  if (focusableElements.length === 0) return
  
  const firstElement = focusableElements[0]
  const lastElement = focusableElements[focusableElements.length - 1]
  
  const handleTabKey = (e: KeyboardEvent) => {
    if (e.key !== 'Tab') return
    
    if (e.shiftKey) {
      // Shift + Tab
      if (document.activeElement === firstElement) {
        e.preventDefault()
        lastElement.focus()
      }
    } else {
      // Tab
      if (document.activeElement === lastElement) {
        e.preventDefault()
        firstElement.focus()
      }
    }
  }
  
  container.addEventListener('keydown', handleTabKey)
  
  // Return cleanup function
  return () => {
    container.removeEventListener('keydown', handleTabKey)
  }
}

/**
 * Generates a unique ID for accessibility purposes
 * @param prefix - Prefix for the ID
 */
export const generateId = (prefix: string = 'id'): string => {
  return `${prefix}-${Math.random().toString(36).substr(2, 9)}`
}

/**
 * Keyboard navigation handler for lists
 * @param event - Keyboard event
 * @param items - Array of focusable items
 * @param currentIndex - Currently focused item index
 * @param onSelect - Callback when item is selected
 */
export const handleListNavigation = (
  event: KeyboardEvent,
  items: HTMLElement[],
  currentIndex: number,
  onSelect?: (index: number) => void
) => {
  const { key } = event
  let newIndex = currentIndex
  
  switch (key) {
    case 'ArrowDown':
      event.preventDefault()
      newIndex = currentIndex < items.length - 1 ? currentIndex + 1 : 0
      break
    case 'ArrowUp':
      event.preventDefault()
      newIndex = currentIndex > 0 ? currentIndex - 1 : items.length - 1
      break
    case 'Home':
      event.preventDefault()
      newIndex = 0
      break
    case 'End':
      event.preventDefault()
      newIndex = items.length - 1
      break
    case 'Enter':
    case ' ':
      event.preventDefault()
      if (onSelect) {
        onSelect(currentIndex)
      }
      return currentIndex
    case 'Escape':
      event.preventDefault()
      ;(event.target as HTMLElement).blur()
      return currentIndex
    default:
      return currentIndex
  }
  
  if (items[newIndex]) {
    items[newIndex].focus()
  }
  
  return newIndex
}

/**
 * Checks if an element is visible and accessible
 * @param element - Element to check
 */
export const isElementAccessible = (element: HTMLElement): boolean => {
  const style = window.getComputedStyle(element)
  
  return (
    style.display !== 'none' &&
    style.visibility !== 'hidden' &&
    style.opacity !== '0' &&
    element.getAttribute('aria-hidden') !== 'true'
  )
}

/**
 * Gets the accessible name of an element
 * @param element - Element to get name for
 */
export const getAccessibleName = (element: HTMLElement): string => {
  // Check aria-labelledby
  const labelledBy = element.getAttribute('aria-labelledby')
  if (labelledBy) {
    const labelElement = document.getElementById(labelledBy)
    if (labelElement) {
      return labelElement.textContent?.trim() || ''
    }
  }
  
  // Check aria-label
  const ariaLabel = element.getAttribute('aria-label')
  if (ariaLabel) {
    return ariaLabel.trim()
  }
  
  // Check associated label element
  if (element.tagName === 'INPUT' || element.tagName === 'TEXTAREA' || element.tagName === 'SELECT') {
    const id = element.getAttribute('id')
    if (id) {
      const label = document.querySelector(`label[for="${id}"]`)
      if (label) {
        return label.textContent?.trim() || ''
      }
    }
  }
  
  // Fallback to text content
  return element.textContent?.trim() || ''
}