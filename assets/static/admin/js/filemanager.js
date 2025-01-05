// static/admin/js/filemanager.js
class FileManager {
    constructor(containerId, options = {}) {
        this.container = document.getElementById(containerId);
        this.currentPath = '/';
        this.options = {
            allowDelete: true,
            allowUpload: true,
            allowRename: true,
            ...options
        };
        this.initialize();
    }

    async initialize() {
        await this.loadCurrentDirectory();
        this.setupEventListeners();
    }

    async loadCurrentDirectory() {
        try {
            const response = await fetch(`/admin/api/files?path=${this.currentPath}`);
            const data = await response.json();
            this.renderFileList(data);
        } catch (error) {
            console.error('Failed to load directory:', error);
        }
    }

    // Add more methods for file operations
}

// Initialize file manager
document.addEventListener('DOMContentLoaded', () => {
    if (document.getElementById('filemanager')) {
        new FileManager('filemanager');
    }
});