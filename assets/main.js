$( document ).ready(function() {
    $( "#buy" ).on( "click", function() {
        var data = {
            "trade_type": "buy",
            "timestamp": $.now().toString(),
            "amount" : $("#amountXTRI").val(),
            "price" : (parseFloat($("#priceXTRI").val()) / 100000000).toFixed(8).toString()
        };
        $.post("/api/buy", data, function(result){
            $("span").html(result);
        });
        
    })
    $( "#sell" ).on( "click", function() {
        var data = {
            "trade_type": "sell",
            "timestamp": $.now().toString(),
            "amount" : $("#amountXTRI").val(),
            "price" : (parseFloat($("#priceXTRI").val()) / 100000000).toFixed(8).toString()
        };
        $.post("/api/sell", data, function(result){
            $("span").html(result);
        });
        
    })
})