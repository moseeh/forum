function timeAgo(date) {
  const now = new Date();
  const diff = now - date;

  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  const months = Math.floor(days / 30);

  if (months >= 1) {
    return date.toLocaleString("en-GB", {
      timeZone: "Africa/Nairobi",
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  } else if (days == 1) {
    return "Yesterday ";
  } else if (days > 1) {
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
  const elements = document.getElementsByClassName("post-date");
  for (let element of elements) {
    const date = new Date(element.textContent.trim() + " GMT");
    element.textContent = timeAgo(date);
  }
}

document.addEventListener("DOMContentLoaded", displayTimeAgo);

function validateForm(event) {
  // Prevent the default form submission
  event.preventDefault();

  // Get form elements
  const form = event.target;
  const title = form.querySelector('input[name="title"]').value;
  const content = form.querySelector('textarea[name="content"]').value;

  // Check if title or content contains only whitespace or is empty
  if (!title.trim()) {
    alert("Title cannot be empty");
    return false;
  }

  if (!content.trim()) {
    alert("Content cannot be empty");
    return false;
  }

  const checkboxes = form.querySelectorAll(".category-checkbox:checked");

  // Validate category selection
  if (checkboxes.length === 0) {
    alert("Please select at least one category");
    return false;
  }

  // If all validations pass, submit the form
  form.submit();
  return true;
}
function validateCommentForm(event) {
  // Prevent the default form submission
  event.preventDefault();
  
  // Get form elements
  const form = event.target;
  const comment = form.querySelector('textarea[name="comment"]').value;
  
  // Check if comment contains only whitespace or is empty
  if (!comment.trim()) {
      alert("Comment cannot be empty");
      return false;
  }
  
  // If validation passes, submit the form
  form.submit();
  return true;
}
