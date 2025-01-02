function timeAgo(date) {
    const now = new Date();
    const diff = now - date;

    const seconds = Math.floor(diff / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);
    const months = Math.floor(days / 30);

    if (months >= 1) {
        return date.toLocaleString('en-GB', {
            timeZone: 'Africa/Nairobi',
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    } else if (days >= 1) {
        return `${days} days ago`;
    } else if (hours >= 1) {
        return `${hours} hours ago`;
    } else if (minutes >= 1) {
        return `${minutes} minutes ago`;
    } else {
        return `${seconds} seconds ago`;
    }
}

function displayTimeAgo() {
    const elements = document.getElementsByClassName('post-date');
    for (let element of elements) {
        const date = new Date(element.textContent.trim() + ' GMT');
        element.textContent = timeAgo(date);
    }
}

document.addEventListener('DOMContentLoaded', displayTimeAgo);