document.addEventListener('DOMContentLoaded', () => {
    const feedbackBtn = document.getElementById('feedbackBtn');
    const modal = document.getElementById('feedbackModalBackdrop');
    const closeBtn = document.getElementById('feedbackModalCloseBtn');
    const form = document.getElementById('feedbackForm');

    if (!feedbackBtn) return;

    feedbackBtn.addEventListener('click', () => modal.classList.add('is-visible'));
    closeBtn.addEventListener('click', () => modal.classList.remove('is-visible'));
    modal.addEventListener('click', e => {
        if (e.target === modal) modal.classList.remove('is-visible');
    });

    form.addEventListener('submit', e => {
        e.preventDefault();
        const type = form.querySelector('input[name="feedbackType"]:checked').value;
        const message = document.getElementById('feedbackMessage').value;
        const page = window.location.pathname;

        fetch('http://localhost:8080/api/feedback', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ type, message, page })
        })
        .then(response => {
            if (!response.ok) throw new Error('Senden fehlgeschlagen');
            return response.json();
        })
        .then(() => {
            modal.classList.remove('is-visible');
            form.reset();
            alert('Danke für dein Feedback!');
        })
        .catch(err => alert(err.message));
    });
});