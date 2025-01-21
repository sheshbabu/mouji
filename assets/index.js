document.addEventListener('DOMContentLoaded', () => {
    initBarChart();
    initDropdown();
});

function initBarChart() {
    const tooltip = document.querySelector('.tooltip');
    const tooltipLabel = document.querySelector('.tooltip > .label');
    const tooltipValue = document.querySelector('.tooltip > .value');
    const bars = document.querySelectorAll('.bar');

    bars.forEach(bar => {
        bar.addEventListener('mousemove', (e) => {
            tooltipLabel.textContent = bar.getAttribute('data-label');
            tooltipValue.textContent = `${bar.getAttribute('data-value')} views`;
            tooltip.style.display = 'flex';

            const barRect = bar.getBoundingClientRect();
            const tooltipRect = tooltip.getBoundingClientRect();

            const tooltipHeight = tooltipRect.height;
            const tooltipWidth = tooltipRect.width;

            const topY = barRect.top + window.scrollY; // https://stackoverflow.com/q/41576287
            let tooltipLeft = barRect.x + (barRect.width / 2) - (tooltipWidth / 2);
            if (tooltipLeft + tooltipWidth > window.innerWidth) { // https://o7planning.org/12397/javascript-window
                tooltipLeft = barRect.right - tooltipWidth;
            } else if (tooltipLeft < 0) {
                tooltipLeft = barRect.left;
            }
            tooltip.style.left = `${tooltipLeft}px`;
            tooltip.style.top = `${topY - tooltipHeight - 5}px`;
        });

        bar.addEventListener('mouseleave', () => {
            tooltip.style.display = 'none';
        });
    });
}

function initDropdown() {
    const dropdownButtons = document.querySelectorAll('.dropdown-button');
    dropdownButtons.forEach(button => {
        button.addEventListener('click', (e) => {
            const dropdown = e.target.closest('.dropdown-container');
            dropdown.classList.toggle('open');
            e.stopPropagation();
        });
    });
}

document.addEventListener('click', (e) => {
    if (e.target.closest('.dropdown-container')) {
        return;
    }

    const dropdowns = document.querySelectorAll('.dropdown-container.open');
    dropdowns.forEach(el => el.classList.remove('open'));
});