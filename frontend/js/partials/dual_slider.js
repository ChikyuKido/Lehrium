document.addEventListener("DOMContentLoaded",function () {
    console.log("test")
    const minSlider = document.getElementById("minSlider");
    const maxSlider = document.getElementById("maxSlider");

    function updateRatingDisplay() {
        const minValue = parseInt(minSlider.value);
        const maxValue = parseInt(maxSlider.value);

        if (minValue > maxValue) {
            minSlider.value = maxValue;
        }
        if (maxValue < minValue) {
            maxSlider.value = minValue;
        }
    }

    minSlider.addEventListener("input", updateRatingDisplay);
    maxSlider.addEventListener("input", updateRatingDisplay);
});