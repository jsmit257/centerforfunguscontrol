$(function () {
  var $lifecycle = $('body>.main>.workspace>.lifecycle')
  var $form = $lifecycle.find('form[name=lifecycle]')

  $form.find('ul>li.attr.expandable').on('click', '.label', e => {
    console.log($(e.target).toggleClass('expanded'))
  })


  $form.on('rx', (e, data) => {
    console.log('message recieved', data)
    var $items = $(e.target).find('>ul.lifecycle')
    $items.find('li.attr[name=id]>.input>.static').text(data.id)
    $items.find('li.attr[name=location]>.input>.static').text(data.location)
    $items.find('li.attr[name=strain_cost]>.input>.static').text(data.strain_cost)
    $items.find('li.attr[name=grain_cost]>.input>.static').text(data.grain_cost)
    $items.find('li.attr[name=bulk_cost]>.input>.static').text(data.bulk_cost)
    $items.find('li.attr[name=yield]>.input>.static').text(data.yield)
    $items.find('li.attr[name=count]>.input>.static').text(data.count)
    $items.find('li.attr[name=gross]>.input>.static').text(data.gross)
    $items.find('li.attr[name=modified_date]>.input>.static').text(data.modified_date)
    $items.find('li.attr[name=create_date]>.input>.static').text(data.create_date)
  })

  $form.on('tx', (e) => {
    var $items = $(e.target).find('>ul.lifecycle')
    $.ajax({
      url: "/lifecycle",
      method: 'POST',
      data: {
        "id": $items.find('li.attr[name=id]>.input>input.live').val(),
        "location": $items.find('li.attr[name=location]>.input>input.live').val(),
        "strain_cost": $items.find('li.attr[name=strain_cost]>.input>input.live').val(),
        "grain_cost": $items.find('li.attr[name=grain_cost]>.input>input.live').val(),
        "bulk_cost": $items.find('li.attr[name=bulk_cost]>.input>input.live').val(),
        "yield": $items.find('li.attr[name=yield]>.input>input.live').val(),
        "count": $items.find('li.attr[name=count]>.input>input.live').val(),
        "gross": $items.find('li.attr[name=gross]>.input>input.live').val(),
        "strain": {
          "id": $items.find('ul[name=strain]>li>.input>input>.live').val()
        },
        "grain_substrate": {
          "id": $items.find('ul[name=grain]>li>.input>input>.live').val()
        },
        "bulk_substrate": {
          "id": $items.find('ul[name=bulk]>li>.input>input>.live').val()
        },
      },
      context: $form,
      async: true,
      success: (result, status, xhr) => {
      },
      error: (xhr, status, err) => {
      },
    })
  })

  $lifecycle.find('.ndx>.rows').on('click', '.row', e => {
    var $row = $(e.currentTarget)
    $.ajax({
      url: "/lifecycle/" + $row.find('div.id').text(),
      method: 'GET',
      context: $form,
      async: true,
      success: (result, status, xhr) => {
        console.log('triggering with result', result)
        $form.trigger('rx', result)
      },
      error: (xhr, status, err) => {
        console.log(xhr, status, err)
      },
    })
    $row.parent().find('.row').removeClass('selected')
    $row.addClass('selected')
  })

  var holybuckets = {
    "id": "0",
    "name": "reference implementation",
    "location": "testing",
    "grain_cost": 1,
    "bulk_cost": 2,
    "yield": 3,
    "count": 4,
    "gross": 5,
    "modified_date": "1970-01-01T00:00:00Z",
    "create_date": "1970-01-01T00:00:00Z",
    "strain": {
      "id": "0",
      "name": "Morel",
      "vendor": {
        "id": "0",
        "name": "127.0.0.1"
      },
      "attributes": [
        {
          "id": "0",
          "name": "contamination resistance",
          "value": "high"
        },
        {
          "id": "1",
          "name": "headroom (cm)",
          "value": "25"
        }
      ]
    },
    "grain_substrate": {
      "id": "0",
      "name": "Rye",
      "type": "Grain",
      "vendor": {
        "id": "0",
        "name": "127.0.0.1"
      },
      "ingredients": [
        {
          "id": "2",
          "name": "Rye"
        }
      ]
    },
    "bulk_substrate": {
      "id": "2",
      "name": "Cedar chips",
      "type": "Bulk",
      "vendor": {
        "id": "0",
        "name": "127.0.0.1"
      }
    },
    "events": [
      {
        "id": "0",
        "temperature": 2,
        "humidity": 1,
        "modified_date": "1970-01-01T00:00:00Z",
        "create_date": "1970-01-01T00:00:00Z",
        "event_type": {
          "id": "1",
          "name": "Fruiting",
          "severity": "Info",
          "stage": {
            "id": "1",
            "name": "Colonization"
          }
        }
      },
      {
        "id": "2",
        "temperature": 0,
        "humidity": 8,
        "modified_date": "1970-01-01T00:00:00Z",
        "create_date": "1970-01-01T00:00:00Z",
        "event_type": {
          "id": "0",
          "name": "Condensation",
          "severity": "Warn",
          "stage": {
            "id": "3",
            "name": "Vacation"
          }
        }
      }
    ]
  }

  $lifecycle.find('.ndx>.rows.row').trigger('click')
})