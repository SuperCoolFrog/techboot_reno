% Test
% parse_command(Any, Action, Arg) :- Action=Any, Arg=42.
%
% A sentence parses a noun phrase, then hands the remaining tokens to the verb phrase
sentence(In, Out) :- 
    noun_phrase(In, Mid), 
    verb_phrase(Mid, Out).

% A noun phrase parses a determiner, then hands the remainder to a noun
noun_phrase(In, Out) :- 
    determiner(In, Mid), 
    noun(Mid, Out).

% A verb phrase parses a verb, then hands the remainder to a noun phrase
verb_phrase(In, Out) :- 
    verb(In, Mid), 
    noun_phrase(Mid, Out).

% Terminal elements unify the head of the list and return the tail (remainder)
determiner([the|Tail], Tail).
determiner([a|Tail], Tail).

noun([cat|Tail], Tail).
noun([dog|Tail], Tail).
noun([mouse|Tail], Tail).

verb([chases|Tail], Tail).
verb([eats|Tail], Tail).

% lex2(In, Out) :- Out = In. 

% lex(Tokens) --> skip_unused, tokens(Tokens).
% 
% tokens([]) --> [].
% 
% tokens([Token|Tokens]) --> 
%     word(Chars), 
%     { Chars \= [] },
%     { atom_chars(Token, Chars) }, 
%     skip_unused, 
%     tokens(Tokens).
% 
% word([C|Cs]) --> [C], { char_type(C, alnum) }, !, word(Cs).
% word([])     --> [].
% 
% skip_unused --> [C], { char_type(C, punct) }, !, skip_unused.
% skip_unused --> [C], { char_type(C, space) }, !, skip_unused.
% skip_unused --> [].

