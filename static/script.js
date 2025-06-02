document.addEventListener('DOMContentLoaded', () => {
    const prevBtn = document.getElementById('prevBtn');
    const nextBtn = document.getElementById('nextBtn');
    const pageInfo = document.getElementById('pageInfo');
    const tbody = document.querySelector('#wordTable tbody');
    
    const { sessionID, currentPage, totalItems, pageSize } = window.initialData;
    let totalPages = Math.ceil(totalItems / pageSize);
    let currentPageNum = currentPage;

    function loadPage(page) {
        if (page < 1 || page > totalPages) return;
        
        fetch(`/data/${sessionID}/${page}`)
            .then(response => response.json())
            .then(data => {
                tbody.innerHTML = '';
                data.words.forEach(word => {
                    const row = document.createElement('tr');
                    row.innerHTML = `
                        <td>${word.Word}</td>
                        <td>${word.TF}</td>
                        <td>${word.IDF.toFixed(4)}</td>
                    `;
                    tbody.appendChild(row);
                });
                
                currentPageNum = page;
                pageInfo.textContent = `Page ${currentPageNum} of ${totalPages}`;
                
                prevBtn.disabled = currentPageNum === 1;
                nextBtn.disabled = currentPageNum === totalPages;
            })
            .catch(error => {
                console.error('Error loading page:', error);
            });
    }

    prevBtn.addEventListener('click', () => {
        loadPage(currentPageNum - 1);
    });

    nextBtn.addEventListener('click', () => {
        loadPage(currentPageNum + 1);
    });

    prevBtn.disabled = currentPageNum === 1;
    nextBtn.disabled = currentPageNum === totalPages;
});