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
});