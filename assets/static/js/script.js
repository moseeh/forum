document.addEventListener('DOMContentLoaded', function () {
    // Sidebar navigation
    document.querySelectorAll('.sidebar .nav-item').forEach(item => {
        item.addEventListener('click', function () {
            // Remove active class from all items
            document.querySelectorAll('.sidebar .nav-item').forEach(i => i.classList.remove('active'));

            // Add active class to the clicked item
            this.classList.add('active');
        });
    });

    // Signup button
    const signupButton = document.getElementById('join-community-btn');
    if (signupButton) {
        signupButton.addEventListener('click', function () {
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


document.addEventListener('DOMContentLoaded', function () {
    // Handle comment button clicks
    document.querySelectorAll('.comment-button').forEach(button => {
        button.addEventListener('click', function () {
            const postId = this.dataset.postId;
            const formContainer = document.getElementById(`comment-form-${postId}`);
            formContainer.style.display = 'block';
        });
    });

    // Handle cancel button clicks
    document.querySelectorAll('.cancel-comment').forEach(button => {
        button.addEventListener('click', function () {
            const formContainer = this.closest('.comment-form-container');
            formContainer.style.display = 'none';
        });
    });
});

document.addEventListener('DOMContentLoaded', function () {
    // Get all necessary elements
    const navItems = document.querySelectorAll('.nav-item');
    const posts = document.querySelectorAll('.post-card');
    const allPosts = document.getElementById('allPosts');
    const likedPosts = document.getElementById('likedPosts');
    const createdPosts = document.getElementById('createdPosts');
    const noPostsMessage = document.getElementById('noPostsMessage');

    navItems.forEach(item => {
        item.addEventListener('click', function (e) {
            e.preventDefault();

            // Remove active class from all items
            navItems.forEach(nav => nav.classList.remove('active'));
            // Add active class to clicked item
            this.classList.add('active');

            const filter = this.dataset.filter;
            console.log('Selected filter:', filter); // Debug log

            // Handle different filter types
            switch (filter) {
                case 'liked':
                    showSection(likedPosts);
                    break;
                case 'all':
                    showSection(allPosts);
                    showAllPosts();
                    break;
                case 'created':
                    showSection(createdPosts);
                    break;
                default:
                    // Category filtering
                    filterByCategory(this.textContent.trim());
            }
        });
    });

    // Helper function to show a section and hide others
    function showSection(sectionToShow) {
        [allPosts, likedPosts, createdPosts].forEach(section => {
            if (section) {
                section.style.display = section === sectionToShow ? 'block' : 'none';
            }
        });
        if (noPostsMessage) {
            noPostsMessage.style.display = 'none';
        }
    }

    // Helper function to show all posts
    function showAllPosts() {
        posts.forEach(post => {
            post.style.display = 'block';
        });
    }

    // Helper function to filter by category
    function filterByCategory(category) {
        showSection(allPosts);
        let visibleCount = 0;

        posts.forEach(post => {
            try {
                const categoryData = post.dataset.category;

                // Debug log for category data
                console.log('Post categories:', {
                    post: post,
                    categoryData: categoryData
                });

                if (categoryData) {
                    const categories = categoryData.split(',')
                        .map(c => c.trim())
                        .filter(Boolean);

                    if (categories.includes(category)) {
                        post.style.display = 'block';
                        visibleCount++;
                    } else {
                        post.style.display = 'none';
                    }
                } else {
                    console.warn('No category data found for post:', post);
                    post.style.display = 'none';
                }
            } catch (error) {
                console.error('Error processing post:', error);
                post.style.display = 'none';
            }
        });

        // Show/hide no posts message
        if (noPostsMessage) {
            noPostsMessage.style.display = visibleCount === 0 ? 'block' : 'none';
        }
    }
});


    document.getElementById('closeModal').onclick = function() {
        document.getElementById('create-post-modal').style.display = 'none';
    };

    // Optional: Close modal on outside click
    window.onclick = function(event) {
        const modal = document.getElementById('create-post-modal');
        if (event.target == modal) {
            modal.style.display = 'none';
        }
    };

