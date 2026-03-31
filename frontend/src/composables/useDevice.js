import { ref, onMounted, onUnmounted } from 'vue'

export function useDevice() {
    const isMobile = ref(false)
    const isTablet = ref(false)
    const isDesktop = ref(true)
    const platform = ref('desktop')

    const updateDevice = () => {
        const width = window.innerWidth
        const ua = navigator.userAgent.toLowerCase()
        
        // Intelligent detection: check for mobile/tablet strings OR small width
        const isMobileUA = /android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i.test(ua)
        
        isMobile.value = width < 768 || (isMobileUA && width < 1024 && !ua.includes('ipad'))
        isTablet.value = (width >= 768 && width < 1024) || (isMobileUA && ua.includes('ipad'))
        isDesktop.value = !isMobile.value && !isTablet.value

        if (isMobile.value) platform.value = 'mobile'
        else if (isTablet.value) platform.value = 'tablet'
        else platform.value = 'desktop'
    }

    onMounted(() => {
        updateDevice()
        window.addEventListener('resize', updateDevice)
    })

    onUnmounted(() => {
        window.removeEventListener('resize', updateDevice)
    })

    return {
        isMobile,
        isTablet,
        isDesktop,
        platform
    }
}
