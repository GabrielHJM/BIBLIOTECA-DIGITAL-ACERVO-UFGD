<template>
  <div class="ios-notification-container">
    <transition-group name="ios-pill" tag="div" class="ios-notification-list">
      <div 
        v-for="notification in notifications" 
        :key="notification.id"
        class="ios-notification-pill"
        :class="['type-' + notification.type]"
        @click="remove(notification.id)"
      >
        <div class="ios-notification-icon">
          <v-icon size="18">{{ getIcon(notification.type) }}</v-icon>
        </div>
        <div class="ios-notification-text">{{ notification.text }}</div>
      </div>
    </transition-group>
  </div>
</template>

<script>
import { iosNotificationStore } from '@/services/IOSNotificationStore'

export default {
  name: 'IOSNotification',
  computed: {
    notifications() {
      return iosNotificationStore.notifications;
    }
  },
  methods: {
    remove(id) {
      iosNotificationStore.removeNotification(id);
    },
    getIcon(type) {
      const icons = {
        success: 'mdi-check-circle',
        error: 'mdi-alert-circle',
        warning: 'mdi-alert',
        info: 'mdi-information',
        system: 'mdi-apple-keyboard-command'
      };
      return icons[type] || icons.info;
    }
  }
}
</script>

<style scoped>
.ios-notification-container {
  position: fixed;
  top: 30px;
  left: 0;
  right: 0;
  display: flex;
  justify-content: center;
  pointer-events: none;
  z-index: 9999;
}

.ios-notification-list {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.ios-notification-pill {
  pointer-events: auto;
  display: flex;
  align-items: center;
  padding: 12px 20px 12px 14px;
  background: rgba(15, 15, 15, 0.35);
  backdrop-filter: blur(40px) saturate(250%);
  -webkit-backdrop-filter: blur(40px) saturate(250%);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 100px;
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.4);
  cursor: pointer;
  max-width: 90vw;
  width: max-content;
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.3s ease;
}

.ios-notification-pill:hover {
  transform: scale(0.98);
}

.ios-notification-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  margin-right: 12px;
  flex-shrink: 0;
}

.ios-notification-text {
  color: #ffffff;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: -0.2px;
  line-height: 1.3;
}

/* Colors for types */
.type-success .ios-notification-icon { color: #32D74B; background: rgba(50, 215, 75, 0.15); }
.type-error .ios-notification-icon { color: #FF453A; background: rgba(255, 69, 58, 0.15); }
.type-warning .ios-notification-icon { color: #FF9F0A; background: rgba(255, 159, 10, 0.15); }
.type-info .ios-notification-icon { color: #0A84FF; background: rgba(10, 132, 255, 0.15); }
.type-system .ios-notification-icon { color: #ffffff; background: rgba(255, 255, 255, 0.15); }

/* Vue Transitions */
.ios-pill-enter-active {
  transition: all 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}
.ios-pill-leave-active {
  transition: all 0.3s cubic-bezier(0.55, 0.085, 0.68, 0.53);
}
.ios-pill-enter-from {
  opacity: 0;
  transform: translateY(-40px) scale(0.8);
}
.ios-pill-leave-to {
  opacity: 0;
  transform: translateY(-20px) scale(0.9);
}
</style>
