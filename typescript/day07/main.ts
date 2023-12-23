import * as fs from 'fs';

enum TypeOfHand {
  HighCard,
  OnePair,
  TwoPairs,
  ThreeOfAKind,
  FullHouse,
  FourOfAKind,
  FiveOfAKind,
}

type Hand = {
  strRepresentation: string
  handType: TypeOfHand
  bid: number
}
function loadData(filename: string): [string, number][] {
  let data = fs.readFileSync(filename, 'utf8').split('\n');
  let result: [string, number][] = [];
  for (let line of data) {
    let [hand, bid] = line.split(' ');
    result.push([hand, parseInt(bid)]);
  }
  return result;
}

function countCards(hand: string): Map<string, number> {
  let result = new Map<string, number>();
  for (let card of hand) {
    if (result.has(card)) {
      result.set(card, result.get(card)! + 1);
    } else {
      result.set(card, 1);
    }
  }
  return result;
}

function parseValues(values: number[]): TypeOfHand {
  if (values[0] == 5) {
    return TypeOfHand.FiveOfAKind;
  } else if (values[0] == 4) {
    return TypeOfHand.FourOfAKind;
  } else if (values[0] == 3 && values[1] == 2) {
    return TypeOfHand.FullHouse;
  } else if (values[0] == 3) {
    return TypeOfHand.ThreeOfAKind;
  } else if (values[0] == 2 && values[1] == 2) {
    return TypeOfHand.TwoPairs;
  } else if (values[0] == 2) {
    return TypeOfHand.OnePair;
  }
  return TypeOfHand.HighCard;
}

function parseHandPart1(hand: string): TypeOfHand {
  let count = countCards(hand);
  let values = Array.from(count.values());
  values.sort((a, b) => b - a);
  return parseValues(values);
}

function parseHandPart2(hand: string): TypeOfHand {
  let count = countCards(hand);
  let jCount = count.get('J') || 0;
  if (jCount == 5 || jCount == 4) {
    return TypeOfHand.FiveOfAKind;
  }
  count.delete('J');
  let values = Array.from(count.values());
  values.sort((a, b) => b - a);
  values[0] += jCount;
  return parseValues(values);
}

function cardValue(card: string, jValue: number): number {
  if (card == 'A') {
    return 14;
  } else if (card == 'K') {
    return 13;
  } else if (card == 'Q') {
    return 12;
  } else if (card == 'J') {
    return jValue;
  } else if (card == 'T') {
    return 10;
  }
  return parseInt(card);
}

function createCmpFunction(jValue: number): (a: Hand, b: Hand) => number {
  return (a: Hand, b: Hand): number => {
    if (a.handType == b.handType) {
      for (let i = 0; i < a.strRepresentation.length; i++) {
        let aCard = a.strRepresentation[i];
        let bCard = b.strRepresentation[i];
        if (aCard == bCard) {
          continue;
        }
        return cardValue(aCard, jValue) - cardValue(bCard, jValue);
      }
    }
    return a.handType - b.handType;
  }
}

function solution(parseFn: (hand: string) => TypeOfHand, jValue: number) {
  let data = loadData('input.txt');
  let hands = data.map(([hand, bid]) => {
    return {
      strRepresentation: hand,
      handType: parseFn(hand),
      bid: bid
    }
  });
  hands.sort(createCmpFunction(jValue));
  let winnings = hands.reduce((prev, curr, index) => {
    return prev + (index +1) * curr.bid;
  }, 0);
  console.log(winnings);
}
function solutionPart1() {
  solution(parseHandPart1, 11);
  solution(parseHandPart2, 1)
}

solutionPart1()