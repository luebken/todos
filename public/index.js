document.addEventListener('DOMContentLoaded', function() {
    const usernameModal = new bootstrap.Modal(document.getElementById('usernameModal'), {
        backdrop: 'static',
        keyboard: false
    });

    const usernameDisplay = document.getElementById('username-display');
    const usernameInput = document.getElementById('username');
    const username = sessionStorage.getItem('username');
    
    if (!username) {
        usernameModal.show();
    } else {
        usernameDisplay.textContent = `List for: ${username}`;
        loadTodos(username);
    }

    document.getElementById('save-username').addEventListener('click', function() {
        const usernameValue = usernameInput.value;
        if (usernameValue) {
            sessionStorage.setItem('username', usernameValue);
            usernameDisplay.textContent = `List for: ${usernameValue}`;
            usernameModal.hide();
            loadTodos(usernameValue);
        } else {
            alert('Please enter a name for this TODO list');
        }
    });

    document.getElementById('todo-form').addEventListener('submit', function(event) {
        event.preventDefault();
        const username = sessionStorage.getItem('username');
        const itemInput = this.querySelector('input[name="Item"]');
        const item = itemInput ? itemInput.value : '';

        console.log('Username:', username);
        console.log('Item:', item);

        if (!item) {
            console.error('Item input not found or empty');
            return;
        }

        console.log('Form submitted:', { Username: username, Item: item });

        fetch('/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ Username: username, Item: item })
        }).then(res => {
            console.log('Response:', res);
            if (res.status == 200) {
                loadTodos(username);
            }
        }).catch(err => {
            console.error('Error:', err);
        });
    });

    function loadTodos(username) {
        fetch(`/?username=${username}`).then(res => res.text()).then(html => {
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');
            const todosSection = doc.querySelector('.todo-list-section');
            document.querySelector('.todo-list-section').innerHTML = todosSection.innerHTML;
            attachEventListeners();
        }).catch(err => {
            console.error('Error loading todos:', err);
        });
    }

    function attachEventListeners() {
        document.querySelectorAll('.fa-pencil').forEach(element => {
            element.addEventListener('click', function() {
                const item = this.getAttribute('data-item');
                updateDb(item);
            });
        });

        document.querySelectorAll('.fa-trash-o').forEach(element => {
            element.addEventListener('click', function() {
                const item = this.getAttribute('data-item');
                removeFromDb(item);
            });
        });
    }

    function removeFromDb(item){
        fetch(`/delete?item=${item}`, {method: "DELETE"}).then(res =>{
            if (res.status == 200){
                const username = sessionStorage.getItem('username');
                loadTodos(username);
            }
        }).catch(err => {
            console.error('Error deleting todo:', err);
        });
    }

    function updateDb(item) {
        let input = document.getElementById(item)
        let newitem = input.value
        fetch(`/update?olditem=${item}&newitem=${newitem}`, {method: "PUT"}).then(res =>{
            if (res.status == 200){
                alert("Database updated")
                const username = sessionStorage.getItem('username');
                loadTodos(username);
            }
        }).catch(err => {
            console.error('Error updating todo:', err);
        });
    }

    attachEventListeners();
});