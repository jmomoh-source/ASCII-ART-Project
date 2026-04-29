document.addEventListener('DOMContentLoaded', () => {
    const textInput    = document.getElementById('text-input');
    const fontSelect   = document.getElementById('font-select');
    const colorSelect  = document.getElementById('color-select');
    const alignSelect  = document.getElementById('align-select');
    const generateBtn  = document.getElementById('generate-btn');
    const outputText   = document.getElementById('output-text');
    const copyBtn      = document.getElementById('copy-btn');
    const downloadBtn  = document.getElementById('download-btn');
    const charCount    = document.getElementById('char-count');
    const themeToggleBtn = document.getElementById('theme-toggle-btn');
    const iconSun        = document.getElementById('theme-icon-sun');
    const iconMoon       = document.getElementById('theme-icon-moon');

    let currentResult = '';

    // Theme initialization
    const savedTheme = localStorage.getItem('ascii-theme');
    const prefersLight = window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches;

    if (savedTheme === 'light' || (!savedTheme && prefersLight)) {
        document.body.classList.add('light-mode');
        iconSun.classList.add('hidden');
        iconMoon.classList.remove('hidden');
    }

    themeToggleBtn.addEventListener('click', () => {
        document.body.classList.toggle('light-mode');
        const isLight = document.body.classList.contains('light-mode');
        
        if (isLight) {
            iconSun.classList.add('hidden');
            iconMoon.classList.remove('hidden');
            localStorage.setItem('ascii-theme', 'light');
        } else {
            iconSun.classList.remove('hidden');
            iconMoon.classList.add('hidden');
            localStorage.setItem('ascii-theme', 'dark');
        }
    });

    // Character counter
    textInput.addEventListener('input', () => {
        const len = textInput.value.length;
        charCount.textContent = `${len} / 550`;
    });

    // Generate on Cmd/Ctrl+Enter
    textInput.addEventListener('keydown', (e) => {
        if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
            generate();
        }
    });

    // Generate button click
    generateBtn.addEventListener('click', generate);

    // Copy to clipboard
    copyBtn.addEventListener('click', () => {
        if (!currentResult) return;
        navigator.clipboard.writeText(currentResult).then(() => {
            showToast('Copied to clipboard!');
        });
    });

    // Download as .txt
    downloadBtn.addEventListener('click', () => {
        if (!currentResult) return;
        const blob = new Blob([currentResult], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'ascii-art.txt';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
        showToast('Downloaded ascii-art.txt');
    });

    async function generate() {
        const text = textInput.value.trim();
        if (!text) {
            textInput.focus();
            return;
        }

        generateBtn.disabled = true;
        generateBtn.classList.add('loading');
        const btnText = generateBtn.querySelector('.btn-text');
        btnText.textContent = 'Generating...';

        try {
            const res = await fetch('/generate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    text:  text,
                    font:  fontSelect.value,
                    color: colorSelect.value,
                    align: alignSelect.value,
                }),
            });

            const data = await res.json();

            if (data.error) {
                outputText.textContent = `Error: ${data.error}`;
                outputText.className = 'output-text';
                currentResult = '';
                copyBtn.disabled = true;
                downloadBtn.disabled = true;
            } else {
                currentResult = data.result;
                outputText.textContent = data.result;

                // Apply color class
                const color = colorSelect.value;
                outputText.className = 'output-text' + (color ? ` color-${color}` : '');

                copyBtn.disabled = false;
                downloadBtn.disabled = false;
            }
        } catch (err) {
            outputText.textContent = `Connection error: Is the server running?`;
            outputText.className = 'output-text';
            currentResult = '';
            copyBtn.disabled = true;
            downloadBtn.disabled = true;
        } finally {
            generateBtn.disabled = false;
            generateBtn.classList.remove('loading');
            btnText.textContent = 'Generate ASCII Art';
        }
    }

    function showToast(message) {
        toast.textContent = message;
        toast.classList.add('show');
        setTimeout(() => toast.classList.remove('show'), 2000);
    }
});
