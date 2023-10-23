$(document).ready(function () {
    function search() {
        const searchText = $("#search-text").val().trim();

        // 显示加载动画
        $("#loading").show();

        // 发送API请求
        $.ajax({
            url: window.location.origin + "/chat",
            method: "POST",
            contentType: "application/json",
            data: JSON.stringify({
                searchText: searchText,
            }),
            success: function (response) {
                // 隐藏动画
                $("#loading").hide();

                // 清空结果
                $("#results-container").empty();

                const card = `
                    <div class="result-card">
                        <pre><code>${response}</code></pre>
                    </div>
                `;
                $("#results-container").append(card);

                // 显示结果
                $("#results").show();
            },
            error: function (error) {
                // 隐藏加载动画
                $("#loading").hide();

                // 显示错误
                M.toast({html: "搜索失败，请重试" + error});
            }
        });
    }

    $("#search-btn").on("click", search);

    $("#search-text").on("keydown", (event) => {
        // 如果按下的是回车键（keyCode 为 13）
        if (event.keyCode === 13) {
            event.preventDefault(); // 阻止默认行为（如表单提交）
            search();
        }
    });

    $(".preset-btn").on("click", function () {
        const presetQuery = $(this).text();
        $("#search-text").val(presetQuery).focus();
    });
});
