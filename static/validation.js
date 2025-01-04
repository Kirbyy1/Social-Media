document.getElementById("register-form")?.addEventListener("submit", function(event) {
    event.preventDefault(); // Prevent form submission
    let isValid = true;

    // Get form inputs
    const username = document.getElementById("username").value.trim();
    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value;
    const confirmPassword = document.getElementById("confirm-password").value;

    // Get error message elements
    const usernameError = document.getElementById("username-error");
    const emailError = document.getElementById("email-error");
    const passwordError = document.getElementById("password-error");
    const confirmPasswordError = document.getElementById("confirm-password-error");

    // Clear previous error messages
    usernameError.textContent = "";
    emailError.textContent = "";
    passwordError.textContent = "";
    confirmPasswordError.textContent = "";

    // Validate username
    if (username.length < 3) {
        usernameError.textContent = "Username must be at least 3 characters.";
        isValid = false;
    }

    // Validate email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
        emailError.textContent = "Please enter a valid email address.";
        isValid = false;
    }

    // Validate password
    if (password.length < 6) {
        passwordError.textContent = "Password must be at least 6 characters.";
        isValid = false;
    }

    // Validate confirm password
    if (password !== confirmPassword) {
        confirmPasswordError.textContent = "Passwords do not match.";
        isValid = false;
    }

    // If valid, submit the form
    if (isValid) {
        alert("Registration successful!");
        this.submit(); // Submit the form
    }
});

document.getElementById("login-form")?.addEventListener("submit", function(event) {
    event.preventDefault(); // Prevent form submission
    let isValid = true;

    // Get form inputs
    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("password").value;

    // Get error message elements
    const usernameError = document.getElementById("username-error");
    const passwordError = document.getElementById("password-error");

    // Clear previous error messages
    usernameError.textContent = "";
    passwordError.textContent = "";

    // Validate username
    if (username.length < 3) {
        usernameError.textContent = "Username must be at least 3 characters.";
        isValid = false;
    }

    // Validate password
    if (password.length < 6) {
        passwordError.textContent = "Password must be at least 6 characters.";
        isValid = false;
    }

    // If valid, submit the form
    if (isValid) {
        alert("Login successful!");
        this.submit(); // Submit the form
    }
});

function setupCapsLockIndicator(inputId, indicatorId) {
    const inputField = document.getElementById(inputId);
    const capsLockIndicator = document.getElementById(indicatorId);

    inputField.addEventListener('keyup', function (event) {
        if (event.getModifierState('CapsLock')) {
            capsLockIndicator.style.display = 'inline';
        } else {
            capsLockIndicator.style.display = 'none';
        }
    });
}



function togglePassword(inputId) {
    const passwordInput = document.getElementById(inputId);
    const icon = document.getElementById(`${inputId}-icon`);

    if (passwordInput.type === "password") {
        passwordInput.type = "text";
        icon.classList.remove("fa-eye");
        icon.classList.add("fa-eye-slash");
    } else {
        passwordInput.type = "password";
        icon.classList.remove("fa-eye-slash");
        icon.classList.add("fa-eye");
    }
}