$( document ).ready(function() {



    $('#tradeTable').DataTable();
    $( "#buy" ).on( "click", function() {
        var seconds = Math.floor(new Date() / 1000);
        var data = {
            "trade_type": "buy",
            "timestamp": seconds.toString(),
            "amount" : $("#amountXTRI").val(),
            "price" : (parseFloat($("#priceXTRI").val()) / 100000000).toFixed(8).toString()
        };
        $.post("/api/buy", data, function(result){
            $("span").html(result);
        });
        
    })
    $( "#sell" ).on( "click", function() {
        var seconds = Math.floor(new Date() / 1000);
        var data = {
            "trade_type": "sell",
            "timestamp": seconds.toString(),
            "amount" : $("#amountXTRI").val(),
            "price" : (parseFloat($("#priceXTRI").val()) / 100000000).toFixed(8).toString()
        };
        $.post("/api/sell", data, function(result){
            $("span").html(result);
        });
        
    })
    setInterval(function(){
        

        $.get("/api/get_trades", function(result){
        $.get("https://blockchain.info/ticker", function(result2){
            var btcusd = 0;

            btcusd = result2["USD"]["15m"];
            var xtriusd = (btcusd * parseFloat(result[0]['Price'])).toFixed(4);
            $("tbody").empty();
            $("#sats").html("Satoshis: " + result[0]['Price'])
            $("#usd").html("USD: " + xtriusd);
            for(var i = 0; i < result.length;i++){
                $("tbody").append(
                    '<tr>' +
                    '<td>' + result[i]["TimeStamp"] + '</td>' +
                    '<td>' + result[i]["TradeType"] + '</td>' +
                    '<td>' + result[i]["Amount"] + ' XTRI</td>' +
                    '<td>' + result[i]["Price"] + ' BTC</td>' +
                    '<td>' + (parseFloat(result[i]["Price"]) * parseFloat(result[i]["Amount"])).toString() + ' BTC </td>' +
                    '</tr>'
                )
            }
        })

        });
        $('#tradeTable').DataTable();
    },15000)

})
