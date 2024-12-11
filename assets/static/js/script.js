document.addEventListener('DOMContentLoaded', function() {
    // Sidebar navigation
    document.querySelectorAll('.sidebar .nav-item').forEach(item => {
        item.addEventListener('click', function() {
            // Remove active class from all items
            document.querySelectorAll('.sidebar .nav-item').forEach(i => i.classList.remove('active'));
            
            // Add active class to the clicked item
            this.classList.add('active');
        });
    });

    // Signup button
    const signupButton = document.getElementById('join-community-btn');
    if (signupButton) {
        signupButton.addEventListener('click', function() {
            window.location.href = '/register';
        });
    }
});

document.addEventListener('DOMContentLoaded', () => {
    const createPostBtn = document.getElementById('create-post-btn');
    const createPostModal = document.getElementById('create-post-modal');
    
    createPostBtn.addEventListener('click', () => {
      createPostModal.style.display = 'block';
    });
  
    // Close modal when clicking outside
    window.addEventListener('click', (e) => {
      if (e.target === createPostModal) {
        createPostModal.style.display = 'none';
      }
    });
  });