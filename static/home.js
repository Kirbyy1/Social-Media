// home.js

document.addEventListener('DOMContentLoaded', function() {
    const postContent = document.getElementById('postContent');
    const postButton = document.getElementById('postButton');

    postContent.addEventListener('input', function() {
        if (postContent.value.trim().length > 0) {
            postButton.disabled = false;
        } else {
            postButton.disabled = true;
        }
    });

    // Function to format time using date-fns
    function formatTime() {
        const timeElements = document.querySelectorAll('.post-time');
        timeElements.forEach(timeElement => {
            const rawTime = timeElement.getAttribute('datetime');
            const parsedTime = new Date(rawTime);
            const now = new Date();

            const differenceInMinutes = dateFns.differenceInMinutes(now, parsedTime);
            const differenceInHours = dateFns.differenceInHours(now, parsedTime);
            const differenceInDays = dateFns.differenceInDays(now, parsedTime);

            let formattedTime;
            if (differenceInMinutes < 60) {
                formattedTime = `${differenceInMinutes}m ago`;
            } else if (differenceInHours < 24) {
                formattedTime = `${differenceInHours}h ago`;
            } else if (differenceInDays < 7) {
                formattedTime = `${differenceInDays}d ago`;
            } else {
                formattedTime = dateFns.format(parsedTime, 'dd/MM');
            }

            timeElement.textContent = formattedTime;
        });
    }

    formatTime();
});