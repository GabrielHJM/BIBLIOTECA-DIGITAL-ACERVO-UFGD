import { reactive } from 'vue';

export const iosNotificationStore = reactive({
    notifications: [],
    
    addNotification(text, type = 'info', timeout = 4000) {
        const id = Date.now() + Math.random().toString(36).substr(2, 9);
        this.notifications.push({ id, text, type });

        if (timeout > 0) {
            setTimeout(() => {
                this.removeNotification(id);
            }, timeout);
        }
        return id;
    },
    
    removeNotification(id) {
        const index = this.notifications.findIndex(n => n.id === id);
        if (index !== -1) {
            this.notifications.splice(index, 1);
        }
    }
});
