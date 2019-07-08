$( document ).ready(function() {
    const chart = LightweightCharts.createChart($("#chart")[0], {
    width: 1000,
    height: 300,
     layout: {
		backgroundColor: '#171B29',
		textColor: 'rgba(255, 255, 255, 0.9)',
	},
	grid: {
		vertLines: {
			color: 'rgba(42, 46, 57, 0)',
		},
		horzLines: {
			color: 'rgba(42, 46, 57, 0.6)',
		},
	},
	crosshair: {
		mode: LightweightCharts.CrosshairMode.Normal,
	},
	priceScale: {
		scaleMargins: {
			top: 0.3,
			bottom: 0.25,
		},
		borderVisible: false,
	},
	timeScale: {
		borderColor: 'rgba(197, 203, 206, 0.8)',
	},
});

var candleSeries = chart.addCandlestickSeries({
    upColor: '#26a69a',
    downColor: '#171B29',
    borderDownColor: '#D1474A',
    borderUpColor: '#26a69a',
    wickDownColor: '#D1474A',
    wickUpColor: '#26a69a',
  });

var volumeSeries = chart.addHistogramSeries({
	color: '#26a69a',
	lineWidth: 2,
	priceFormat: {
		type: 'volume',
	},
	overlay: true,
	scaleMargins: {
		top: 0.8,
		bottom: 0,
	},
});

candleSeries.setData([
	{ time: '2019-07-07', open: 192.54, high: 193.86, low: 190.41, close: 193.59 },
]);
volumeSeries.setData([
	{ time: '2019-07-07', value: 11487448.00, color: 'rgba(255,82,82, 0.8)' },
	{ time: '2019-07-07', value: 11707083.00, color: 'rgba(255,82,82, 0.8)' },
	{ time: '2019-07-07', value: 8755506.00, color: 'rgba(0, 150, 136, 0.8)' },
	{ time: '2019-07-07', value: 3097125.00, color: 'rgba(0, 150, 136, 0.8)' },
]);
    url = 'ws://127.0.0.1:80/ws';
    c = new WebSocket(url);
    
    send = function(data){
      $("#output").append((new Date())+ " ==> "+data+"\n")
      c.send(data)
    }

    c.onmessage = function(msg){
      $("#output").append((new Date())+ " <== "+msg.data+"\n")
      console.log(msg)
    }

    c.onopen = function(){
      setInterval( 
        function(){ send("heartbeat") }
      , 25000 )
    }

    $( "#buy" ).on( "click", function() {
        var seconds = Math.floor(new Date() / 1000);
        var data = '{ "trade_type" : "buy" , "timestamp": "' + seconds.toString() + '" ,"amount": "' +  $("#amountXTRI").val() + '" ,"price":"' +  (parseFloat($("#priceXTRI").val()) / 100000000).toFixed(8).toString() + '"}';
        console.log(JSON.stringify(data))
        c.send(data)
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
})
