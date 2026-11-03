async function sendForm(form, method = 'POST') {
    const formData = new FormData(form);

    const response = await fetch(form.action, {
        method: method ? method : form.method,
        body: formData
    });

    if (!response.ok) throw new Error('Server error');

    try {
        return await response.json();
    } catch (error) {
        throw new Error('Data error');
    }
}

async function send(url, method = 'GET', data = '', responseJson = true) {
    const response = await fetch(url, {
        method: method,
        data: data
    });

    if (!response.ok) throw new Error('Server error')

    if (responseJson) {
        try {
            return await response.json();
        } catch (error) {
            throw new Error('Data error');
        }
    }
}

function errorHandler(error) {
    console.log(error);
}

async function clientCreate(form) {
    let data;

    try {
        data = await sendForm(form);
    } catch (error) {
        errorHandler(error);
        return;
    }

    afterClientCreate(data);
}

function afterClientCreate(data) {
    location.href = `/clients/${data.uid}`;
}

async function clientUpdate(form) {
    let data;

    try {
        data = await sendForm(form, 'PATCH');
    } catch (error) {
        errorHandler(error);
        return;
    }

    afterClientUpdate(data);
}

function afterClientUpdate(data) {
    location.href = `/clients/${data.uid}`;
}

async function clientDelete(url) {
    try {
        await send(url, 'DELETE', '', false);
    } catch (error) {
        errorHandler(error)
        return;
    }

    afterClientDelete();
}

function afterClientDelete() {
    location.href = '/clients';
}